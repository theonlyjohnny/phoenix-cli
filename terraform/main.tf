provider aws {
  # Env vars:
  # AWS_SECRET_KEY
  # AWS_ACCESS_KEY
  # AWS_DEFAULT_REGION
  # or through credentials file @ ~/.aws/credentials
}

# Find most recent Bionic Ubuntu AMI
data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical LTD.
}
