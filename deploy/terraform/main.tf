# Configure Docker provider and connect to the local Docker socket
provider "docker" {
  host = "unix:///var/run/docker.sock"
}

# Create an evm-lite containers
resource "docker_container" "evm-lite" {
  count = "${var.nodes}"
  
  name = "node${count.index}"
  hostname= "node${count.index}"

  image = "mosaicnetworks/evm-lite:0.1.0"

  networks = ["${docker_network.private_network.name}"]

  volumes {
    host_path = "${path.cwd}/../${var.command}/conf/node${count.index}"
    container_path = "/.evm-lite"
    read_only = true
  }

  # entrypoint =  ["tail", "-f", "/dev/null"]
  command = ["${var.command}"]

  provisioner "local-exec" {
    command = "echo node${count.index} ${self.ip_address}  >> ips.dat"
  }
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

output "public_addresses" {
    value = ["${docker_container.evm-lite.*.ip_address}"]
}