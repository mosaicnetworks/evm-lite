# E2E : End-To-End Tests

This folder contains a few scripts to test and demonstrate certain scenarios on
a local evm-lite/babble network. This includes launching a pre-configured 
testnet with docker containers, adding nodes dynamicaly, and testing a 
crowd-funding smart-contract.

## Dependencies

### Terraform 

We use [Terraform](https://www.terraform.io/) to configure and automate the 
deployment of evm-lite testnets (cf terraform/ directory). To install Terraform: 

```bash
# Download stable release
$ cd /tmp
$ wget https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_amd64.zip

# Unzip and move to local bin 
$ unzip terraform_0.11.13_linux_amd64.zip
$ sudo mv terraform /usr/local/bin/

# Check terraform is available from your standard path. 
$ terraform --version
```
### Node

We use [Node.js](https://nodejs.org) to run the interractive crowd-funding demo,
which is partly written in javascript. The best way to install Node.js is with 
[Node Version Manager](https://github.com/creationix/nvm): 

```bash
# install node version manager 
$ curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.5/install.sh | bash
# use nvm to install and set as default the LTS version of node
$ nvm install node --lts=dubnium
$ nvm alias default lts/dubnium
$ nvm use lts/dubnium
# install dependencies for this demo
evm-lite/e2e$ make deps
```

### evm-lite-cli

`evm-lite-cli` is a command-line client for an `evm-lite` node. We use it in 
some scripts to sign and send transactions.

```bash
# TODO: Installation steps - the line below installs the published release, which may not be new enough.

$ npm install evm-lite-cli -g
```

## Launch a testnet

Create the testnet configuration, and launch the first 4 nodes:

```bash
$ make deploy
```
This loads and starts a prebuilt testnet with 10 nodes, with the first 4 nodes preset as validators. The other 6 nodes are shut down. The last lines of the output shows which nodes are running, and their IP addresses, such as below:

```
node0 172.77.5.3
node1 172.77.5.9
node2 172.77.5.7
node3 172.77.5.8

```

For more information, refer to `/deploy/prebuilt/templates/babblepoademo`, or 
`scripts/`.

## Monitor the running nodes

```bash
$ make watch
```

## Run the crowd-funding demo

```bash
$ make crowd-funding-demo
```
## Add a node

```bash
$ make add-node NODE=node4
```

This combines a few steps:

1) point `evmlc` to the keystore which contains node0's key (so it can sign 
    transactions from node0)

2) node0 nominates node4

3) all nodes vote 'yes' for node4

4) node4 retrieves the current peers.json

5) node4 is started

If everything goes well, node4 should detect that it doesn't belong to the 
Babble validator-set, and proceed to send a join request. The join request goes 
through consensus, and is verified against the POA smart-contract. In prior 
steps (1-3), all the current validators voted to add node4 to the whitelist, so 
node4's request is accepted. Once it is added to the Babble validator-set, node4
joins the gossip routines and syncs with the other nodes.

For a more fine-grained control over the process, please refer to the `scripts/`
folder directly.

## Destroy the testnet

```bash
$ make stop
```