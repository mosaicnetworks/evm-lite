//This file defines the variables used in the main terraform file. Some of the
//private values are specified in a separate file (ex. secret.tfvars) which is
//not included in source control for obvious security reasons.

//AWS API access key
variable "access_key" {}
variable "secret_key" {}

//ID of Virtual Private Cloud
variable "vpc" {}

//PEM file containing RSA key to SSH into AWS instances
variable "key_name" {}

variable "key_path" {}

//evm-lite AMI ID
variable "ami" {
  default = "ami-0691a4aeecbf406ef"
}

//Number of nodes to deploy
variable "nodes" {
  default = 4
}

variable "command" {
  default = "solo"
}

variable "conf" {
  default = ""
}
