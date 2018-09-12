# Configure Docker provider and connect to the local Docker socket
provider "docker" {
  host = "unix:///var/run/docker.sock"
}

# Create an evm-lite container
resource "docker_container" "evm-lite" {
  count = "${var.servers}"
  
  name = "node${count.index}"
  
  image = "mosaicnetworks/evm-lite:0.1.0"

  networks = ["${docker_network.private_network.name}"]
  publish_all_ports = true

  provisioner "file" {
      source = "../conf/node${count.index}/"
      destination = "/.evm-lite"
      connection {
        type     = "ssh"
        agent = false
        host =  "${self.ip_address}"
        private_key = "${file("${path.cwd}/../docker/keys/client")}"
      }
  }

  command = ["solo"]
}

# Create a new docker network
resource "docker_network" "private_network" {
  name = "monet"
  check_duplicate = true
  driver = "bridge"
  ipam_config {
      subnet = "172.77.5.0/24"
      gateway = "172.77.5.254"
  }
}


