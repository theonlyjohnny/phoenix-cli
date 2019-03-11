# data "template_file" "install_salt" {
    # template = "${file("../../terraform/scripts/install_salt.sh")}"
    # vars {
        # salt_version = "${var.salt_version}"
    # }
# }

resource "aws_instance" "phoenix" {

  # The connection block tells our provisioner how to
  # communicate with the resource (instance)
  connection {
    # The default username for our AMI
    user        = "phoenix"
    private_key = "${file(var.phoenix_key_filepath)}"

    # The connection will use the local SSH agent for authentication.
  }

  instance_type = "${var.instance_size}"

  root_block_device {
    volume_size = "500" #TODO this is prolly too big
  }

  instance_initiated_shutdown_behavior = "terminate"

  # Lookup the correct AMI based on the region
  # we specified
  ami           = "${data.aws_ami.ubuntu.id}"

  # The name of our SSH keypair we created above.
  key_name = "${var.phoenix_key_name}"

  # Our Security group to allow HTTP and SSH access
  vpc_security_group_ids = ["${aws_default_security_group.default.id}",
    "${aws_security_group.phoenix.id}"]
  ]


  #Launch into public subnet so we can further configure -- Phoenix master instance should run on private subnet
  subnet_id = "${aws_subnet.subnet_public.0.id}"

  iam_instance_profile = "phoenix"

  tags {
    Name = "${var.vpc_name}-phoenix.${var.aws_domain}"
  }

  # provisioner "file" {
    # source      = "${var.key_filename}"
    # destination = "${var.key_name}"
  # }

  # provisioner "file" {
    # source      = ""
    # destination = "/tmp/pillar_top.sls"
  # }

  # provisioner "file" {
    # destination = "install_salt.sh"
  # }

  provisioner "remote-exec" {
    inline = [
      "set -x",
      "until [ -f /var/lib/cloud/instance/boot-finished ]; do sleep 1 && echo sleep; done",
      "sudo chown 0600 /home/admin/.ssh/id_rsa",
      "sudo chown 0600 /home/admin/.ssh/admin.aws",
      "sudo chown 0600 /home/admin/.ssh/known_hosts",
      "sudo mkdir -p /srv/pillar/base",
      "sudo cp /tmp/pillar_top.sls /srv/pillar/top.sls",
      "sudo rm -rf /srv/salt",
      "sudo mkdir -p /srv/www/aws-saltstack",
      "sudo chown admin -R /srv/www",

      # "curl install_phoenix.sh and run"
    ]
  }
}

resource "aws_eip" "phoenix" {
  vpc      = true
  instance = "${aws_instance.admin.id}"
}
