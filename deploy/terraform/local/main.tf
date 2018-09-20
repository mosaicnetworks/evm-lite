# Configure Docker provider and connect to the local Docker socket
provider "docker" {
  host = "unix:///var/run/docker.sock"
}

# Create a new docker network
resource "docker_network" "private_network" {
  name            = "monet"
  check_duplicate = true
  driver          = "bridge"

  ipam_config {
    subnet  = "172.77.5.0/24"
    gateway = "172.77.5.254"
  }
}

# Create evm-lite containers
resource "docker_container" "evm-lite" {
  count = "${var.nodes}"

  name     = "node${count.index}"
  hostname = "node${count.index}"

  image = "mosaicnetworks/evm-lite:${var.version}"

  networks = ["${docker_network.private_network.name}"]

  volumes {
    host_path      = "${var.conf}/node${count.index}"
    container_path = "/.evm-lite"
    read_only      = false
  }

  # entrypoint = ["tail", "-f", "/dev/null"]
  command = ["${var.command}"]

  provisioner "local-exec" {
    command = "echo node${count.index} ${self.ip_address}>> ips.dat"
  }
}

output "public_addresses" {
  value = ["${docker_container.evm-lite.*.ip_address}"]
}
