##### Setup #####

# Create a VPC to launch our instances into
resource "aws_vpc" "main" {
  cidr_block           = "${var.vpc_cidr_prefix}.0.0"
  enable_dns_support   = true
  enable_dns_hostnames = true

  tags = {
    Name     = "${var.vpc_name}"
  }
}

# DHCP options for the VPC
resource "aws_vpc_dhcp_options" "dhcp_opts" {
  domain_name         = "${var.aws_domain}"
  domain_name_servers = ["AmazonProvidedDNS"]

  tags {
    Name = "${var.vpc_name} DHCP"
  }
}

#Associate DHCP options to VPC
resource "aws_vpc_dhcp_options_association" "vpc_dhcp_options_association" {
  vpc_id          = "${aws_vpc.main.id}"
  dhcp_options_id = "${aws_vpc_dhcp_options.dhcp_opts.id}"
}

# Create a way out to the internet
resource "aws_internet_gateway" "gw" {
  vpc_id = "${aws_vpc.main.id}"

  tags {
    Name = "${var.vpc_name} IGW"
  }
}

# Public route as way out to the internet
resource "aws_route" "internet_access" {
  route_table_id         = "${aws_vpc.main.main_route_table_id}"
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = "${aws_internet_gateway.gw.id}"
}

##### Private #####

# Create the private route table for each availability zone
resource "aws_route_table" "private_route_table" {
  count  = "${len(var.zones)}"
  vpc_id = "${aws_vpc.main.id}"

  tags {
    Name = "${var.vpc_name} Private route table #${count.index}"
  }
}

# Create private route connecting each availability zone
resource "aws_route" "private_route" {
  count                  = "${len(var.zones)}"
  route_table_id         = "${element(aws_route_table.private_route_table.*.id, count.index)}"
  destination_cidr_block = "0.0.0.0/0"
  nat_gateway_id         = "${element(aws_nat_gateway.nat-gateway.*.id, count.index)}"
}

# Create Private subnet (Phoenix instance will launch into one of these)
resource "aws_subnet" "subnet_private" {
  count      = "${len(var.zones)}"
  vpc_id     = "${aws_vpc.main.id}"
  cidr_block = "${var.vpc_cidr_prefix}.${count.index + 1}.0/24"

  availability_zone = "${var.zones[count.index]}"

  tags = {
    Name = "${var.vpc_name} Private subnet #${count.index}"
  }
}

# Associate subnets to private route table
resource "aws_route_table_association" "subnet_private_association" {
  count          = "${len(var.zones)}"
  subnet_id      = "${element(aws_subnet.subnet_private.*.id, count.index)}"
  route_table_id = "${element(aws_route_table.private_route_table.*.id, count.index)}"
}

##### Public #####

# Public subnets
resource "aws_subnet" "subnet_public" {
  count                   = "${len(var.zones)}"
  vpc_id                  = "${aws_vpc.main.id}"
  cidr_block              = "${var.vpc_cidr_prefix}.${10 + count.index * 4}.0/22"
  map_public_ip_on_launch = true
  availability_zone       = "${var.zones[count.index]}"

  tags = {
    Name = "${var.vpc_name} Public subnet #${count.index}"
  }
}

# Associate subnet to public route table
resource "aws_route_table_association" "subnet_public-association" {
  count          = "${len(var.zones)}"
  subnet_id      = "${element(aws_subnet.subnet_public.*.id, count.index)}"
  route_table_id = "${aws_vpc.main.main_route_table_id}"
}

# Create a nat gateway for each availability zone and it will depend on the internet gateway creation
resource "aws_nat_gateway" "nat-gateway" {
  count         = "${len(var.zones)}"
  allocation_id = "${element(aws_eip.nat.*.id, count.index)}"
  subnet_id     = "${element(aws_subnet.subnet_public.*.id, count.index)}"

  tags = {
    Name = "${var.vpc_name} Public NAT Gateway for ${element(aws_subnet.subnet_public.*.id, count.index)}"
  }
}

# Create an EIP for the NAT Gateway in each availability zone
resource "aws_eip" "nat" {
  count      = "${len(var.zones)}"
  vpc        = true
}

##### Security Groups #####

# default security group, mimics default aws default
resource "aws_default_security_group" "default" {
  vpc_id = "${aws_vpc.main.id}"

  ingress {
    protocol  = -1
    self      = true
    from_port = 0
    to_port   = 0
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "${var.vpc_name} Default Security Group"
  }
}

# phoenix setup security group -- allows SSH from box creating instance for further configuration
resource "aws_security_group" "phoenix" {
  name        = "phoenix"
  description = "Phoenix SSH Configuration"
  vpc_id      = "${aws_vpc.main.id}"

  ingress {
    from_port         = 22
    to_port           = 22
    cidr_blocks       = ["${var.self_ip}"]
    protocol          = "tcp"
  }

  tags {
    Name = "${var.vpc_name} Phoenix"
  }
}
