provider "aws" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
  region     = "eu-west-2"
}

resource "aws_subnet" "monet" {
  vpc_id                  = "${var.vpc}"
  cidr_block              = "10.0.2.0/24"
  map_public_ip_on_launch = "true"

  tags {
    Name = "Testnet"
  }
}

resource "aws_security_group" "monetsec" {
  name        = "monetsec"
  description = "MONET internal traffic + maintenance."

  vpc_id = "${var.vpc}"

  // These are for internal traffic
  ingress {
    from_port = 0
    to_port   = 65535
    protocol  = "tcp"
    self      = true
  }

  ingress {
    from_port = 0
    to_port   = 65535
    protocol  = "udp"
    self      = true
  }

  // These are for maintenance
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 1337
    to_port     = 1337
    protocol    = "tcp"
    cidr_blocks = ["10.0.2.0/24"]
  }

  ingress {
    from_port   = -1
    to_port     = -1
    protocol    = "icmp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  // This is for outbound internet access
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "node" {
  count = "${var.nodes}"

  //EVM-LITE v0.1.0 (cf. packer/ to create the AMI)
  ami           = "${var.ami}"
  instance_type = "t2.micro"

  subnet_id              = "${aws_subnet.monet.id}"
  vpc_security_group_ids = ["${aws_security_group.monetsec.id}"]
  private_ip             = "10.0.2.${10+count.index}"

  key_name = "${var.key_name}"

  connection {
    user        = "ubuntu"
    private_key = "${file("${var.key_path}")}"
  }

  provisioner "file" {
    source      = "${var.conf}/node${count.index}"
    destination = "/home/ubuntu/.evm-lite" 
  }

  provisioner "local-exec" {
    command = "echo ${self.private_ip} ${self.public_ip}  >> ips.dat"
  }

  provisioner "remote-exec" {
    inline = [
      "nohup evml ${var.command} > babble_logs 2>&1 &",
      "sleep 1"
      ]
  }

  #Instance tags
  tags {
    Name = "node${count.index}"
  }
}

output "public_addresses" {
  value = ["${aws_instance.node.*.public_ip}"]
}
