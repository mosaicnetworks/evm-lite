# TERRAFORM

This folder contains a set of scripts and configuration files to automate the deployment of evm-lite testnets using Terraform. Tesnets can be deployed locally with Docker containers, or in AWS using using virtual servers. The machine
images for Docker and AWS can be built from the `deploy/build/` folder.

## LOCAL

Local testnets are formed of multiple Docker containers on the same machine. 
They are good to quickly test and play around with evm-lite. 

On Linux, Docker containers are directly accessible from the host machine, so 
you can bootstrap a testnet, and interract with it directly from a separate terminal. On other operating systems, an additional layer of abstraction makes 
it necessary to interract with testnet containers from other containers within the same subnet. 

The scripts will first create a local virtual bridge network called `monet`, where container IPs will be in the `172.77.5.0/24` range (from 172.77.5.0 to 
172.77.5.255). Containers connected to this network will automatically expose 
all ports to each other, and no ports to the outside world. Special ports (for 
the evm-lite HTTP service for example) may be opened from the Dockerfile (cf 
depoly/build/docker) or Terraform main.tf.

Containes are assigned names and hostnames of the form `node0...node4...nodeN`, and can use those hostnames directly to communicate with one-another within the `monet` subnet. To access a container from the host, use the `172.77.5.X` 
address. 

The Docker containers, built from the Dockerfile in deploy/build/docker, come 
pre-packaged with `evml`. Configuration files are mounted through a volume 
attached to the default `/.evm-lite` directory, which is the default location 
for `evml`.  

## AWS

It is also possible to automate the deployment of testnets on AWS. This will create and provision multiple virtual servers in the Amazon Cloud where they can stay running indefinitely and accessible on the public internet. It obviously requires an AWS account and corresponding access keys. Also be aware that deploying resources on AWS in not necessarily free!

There are two types of credentials to provide to Terraform:

- The AWS API Access Key to connect to AWS and provision resources
- An SSH key to communicate with the provisioned instances

These credetials must be created from the AWS console before using these 
scripts. Once created and retrieved from AWS, the credetials must be provided in the `/aws/secret.tfvars` file:

```
//AWS API ACCESS KEY
access_key = "..."
secret_key = "..."

//RSA KEY FOR SSH
key_name = "..."
key_path = "..."
```

The scripts will create an AWS subnet in the `10.0.2.0/24` range and assign it a 
security group defining which ports should remain open or closed for machines
connected to this network. Then it will create a certain number of instances,
built using the evm-lite AMI (cf. deploy/build), and connect them withing this
subnetwork.





