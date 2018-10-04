# DEPLOY

**THIS IS NOT INTENDED FOR PRODUCTION DEPLOYMENTS, JUST TESTING**

We provide a set of scripts to automate the deployment of evm-lite networks 
locally (using docker) or in the cloud (using AWS). Parameters control the 
number of nodes, and which consensus system to use (solo, babble, raft, etc.).

Whether locally or in the cloud, the workflow is as follows:

1. **Build**: Create virtual machine image (docker or ami)
2. **Config**: Generate configuration files for each node and for the network as 
               a whole
3. **Deploy**: Create and start node instances

**The scripts have only been tested on Ubuntu**

## BUILD

### Docker Image

```bash
$ make build-docker-image
```
Produce a versioned Docker image with `evml` using a classic Dockerfile. 

**Dependencies**: Requires Docker Engine and go dependencies (use glide intall 
                  in root dir)

### Amazon Machine Image (AMI)

```bash
$ make build-ami
```
Produce an AMI to run instances (virtual servers) in the AWS cloud with `evml` 
pre-installed.

**Dependencies:** Requires Hashicorp Packer, an AWS account, and AWS access keys

From Packer documentation:

> This builder builds an AMI by launching an EC2 instance from a source AMI, 
> provisioning that running machine, and then creating an AMI from that machine. 
> This is all done in your own AWS account. The builder will create temporary 
> keypairs, security group rules, etc. that provide it temporary access to the 
> instance while the image is being created. This simplifies configuration quite 
> a bit.
> 
> The builder does not manage AMIs. Once it creates an AMI and stores it in your 
> account, it is up to you to use, delete, etc. the AMI.

Provide the AWS access key in the `build/ami/secret.json` file: 

```json
{
    "aws_access_key" : "...",
    "aws_secret_key" : "..."
}
```

## CONFIG

```bash
$ make conf CONSENSUS=[solo] NODES=[1] IPBASE=[node] IPADD=[0]
```

Create the configuration files for the network.

Parameters:

- CONSENSUS: solo, babble, or raft

- NODES: number of nodes in the network

- IPBASE/IPADD: used to determine the address of nodes.
    
ex: If IPBASE=10.0.2. IPADD=10, and NODES=4, the resulting addresses will be:
    10.0.2.10, 10.0.2.11, 10.0.2.12, and 10.0.2.13.
    If IPBASE and IPADD are not specified, the resulting addresses will default
    to: node0, node1, node2, and node3.

The configuration is written to the `conf/[consensus]/conf` folder. For each 
node, there will usually be two configuration sub-directories: one for the 
Ethereum-related things, and one for the consensus related things. For example,
creating the configuration for two babble nodes yields the following files:

`$ make consensus=babble nodes=2`

```bash
conf/babble/conf/
├── evml.toml
├── genesis.json
├── keystore
│   ├── UTC--2018-09-28T15-17-47.839800703Z--370728a25a925a4db0f59ddcd82a0874225bc43b
│   └── UTC--2018-09-28T15-17-49.835572440Z--2c2e382b974716daa13a919fc3143f0cdcef08ce
├── node0
│   ├── babble
│   │   ├── peers.json
│   │   └── priv_key.pem
│   ├── evml.toml
│   └── eth
│       ├── genesis.json
│       ├── keystore
│       │   └── UTC--2018-09-28T15-17-47.839800703Z--370728a25a925a4db0f59ddcd82a0874225bc43b
│       └── pwd.txt
├── node1
│   ├── babble
│   │   ├── peers.json
│   │   └── priv_key.pem
│   ├── evml.toml
│   └── eth
│       ├── genesis.json
│       ├── keystore
│       │   └── UTC--2018-09-28T15-17-49.835572440Z--2c2e382b974716daa13a919fc3143f0cdcef08ce
│       └── pwd.txt
└── peers.json
```

It creats an Ethereum key for each node using the derault password file, and a
config.json file. The genesis file is used by evm-lite to initialize the state 
with some funded accounts. It also created a Babble key and peers.json file
defining the Babble network. The evml.toml file contains parameters for babble
and evm-lite.

## DEPLOY

### Local

Local testnets are formed of multiple Docker containers running on the host 
machine; they are convenient to quickly test evm-lite. 

On Linux, Docker containers are directly accessible from the host, so one can 
bootstrap a testnet, and interact with it directly from a separate terminal. On 
other operating systems, an additional layer of abstraction makes it necessary 
to interact with testnet containers from other containers within the same 
subnet. 

The scripts will first create a local virtual bridge network called `monet`, 
where container IPs will be in the `172.77.5.0/24` range (from 172.77.5.0 to 
172.77.5.255). Containers connected to this network will automatically expose 
all ports to each other, and no ports to the outside world. Special ports (for 
the evm-lite HTTP service for example) may be opened from the Dockerfile (cf 
depoly/build/docker) or Terraform main.tf.

Containers are assigned names and hostnames of the form `node0...node4...nodeN`, 
and can use those hostnames directly to communicate with one-another within the 
`monet` subnet. To access a container from the host, use the `172.77.5.X` 
address. 

The Docker containers, built from the Dockerfile in deploy/build/docker, come 
pre-packaged with `evml`. Configuration files are mounted through a volume 
attached to the default `~/.evm-lite` directory, which is the default location 
for `evml`.  

MACOS USERS

The configuration folders are mounted as volumes in the docker containers, and 
require read/write access from the process running in the container. Hence, when
running the docker containers (cf examples below), specify a user that has 
read/write access to the configuration folders. In most cases, this will be the 
current user (1000 on Linux, 501 on MACOS). 

Examples:

First, build the evm-lite docker image (cf BUILD).
 
``` bash
cd deploy
# configure and start a testnet of 4 evm-lite nodes with Babble consensus
make CONSENSUS=babble NODES=4
# same for MACOS users
make USER=501 CONSENSUS=babble NODES=4 
#configure and start a single evm-lite instance with Solo consensus 
make CONSENSUS=solo NODES=1 
#configure and start a testnet of 3 evm-lite nodes with Raft consensus
make CONSENSUS=raft NODES=3
#bring everything down
make stop 
```

### Cloud

It is also possible to automate the deployment of testnets on AWS. This will 
create and provision multiple virtual servers in the Amazon Cloud where they can 
stay running indefinitely and accessible on the public internet. It obviously 
requires an AWS account and corresponding access keys. Also be aware that 
deploying resources on AWS in not necessarily free!

There are two types of credentials to provide to Terraform:

- The AWS API Access Key to connect to AWS and provision resources
- An SSH key to communicate with the provisioned instances

These credentials must be created from the AWS console before using these 
scripts. Once created and retrieved from AWS, the credetials must be provided in 
the `/aws/secret.tfvars` file:

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
connected to this network. Then it will create a number of instances, built 
using the evm-lite AMI (cf. BUILD), and connect them withing this subnetwork.

Examples:

First, build the evm-lite AMI (cf BUILD) and record its ID in the `ami` 
terraform variable (aws/variables.tf).

```bash
# configure and start a testnet of 4 nodes in AWS
make ENV=aws CONSENSUS=babble NODES=4 IPBASE=10.0.2. IPADD=10
# bring everything down
make stop ENV=aws 
```
