# DEMO

Here we provide a few scripts to demonstrate how to interact with evm-lite
nodes.

You might need to install [https://github.com/creationix/nvm](Node Version Manager) and [https://nodejs.org](NodeJS) and dependencies first:

```bash
# install node version manager 
$ curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.5/install.sh | bash
# use nvm to install and set as default the LTS version of node
$ nvm install node --lts=dubnium
$ nvm alias default lts/dubnium
$ nvm use lts/dubnium
# install dependencies for this demo
evm-lite/demo$ npm install
```

You may also need to install [https://www.terraform.io/](terraform). 
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

Within the ...evm-lite/deploy/terraform/aws directory run:

```bash
terraform init
```

To start a testnet, execute the deploy commands from the `deploy/` directory - the dependencies are described in the [https://github.com/mosaicnetworks/evm-lite/blob/master/deploy/README.md](README) file in the deploy folder:

ex:

```bash
evm-lite$ cd deploy
evm-lite/deploy$ make CONSENSUS=babble NODES=4
```

Then, in an other terminal, start the interactive demo:

```bash
$ ./demo.sh ../deploy/terraform/local/ips.dat
```

The ips.dat file, generated during the deploy phase, tells the demo program
where to reach the nodes.

In this case, we are using Babble consensus, so it is interesting to monitor
the babble nodes:

```bash
$ ./watch.sh ../deploy/terraform/local/ips.dat
```

After the demo, destroy the testnet by running `make stop` from `deploy/`:

```bash
$ cd ../deploy
$ make stop
```
