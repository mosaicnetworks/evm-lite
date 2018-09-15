# DEMO

Here we provide a few scripts to demonstrate how to interract with a testnet of
4 evm-lite nodes using Babble consensus.

To start the testnet, execute the deploy commands from the `deploy/` directory:

```bash
evm-lite$ cd deploy
evm-lite/deploy$ make consensus=babble nodes=4
```

Then, from this directory, launch the `watch` script, which monitors the Babble
status of all 4 nodes:

```bash
evm-lite/deploy$ cd ../demo
evm-lite/demo$ ./watch
```

And in an other terminal, start the interractive demo:

```bash
evm-lite/demo$ ./demo.sh
```

You might need to install NodeJS and dependencies first:

```bash
# install node version manager
evm-lite/demo$ curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.5/install.sh | bash
# use nvm to intall stable version of node
evm-lite/demo$ nvm install node stable
# install dependencies for this demo
evm-lite/demo$ npm install
```

After the demo, destroy the testnet by running `make stop` from `deploy/`:

```bash
evm-lite/demo$ cd ../deploy
evm-lite/deploy$ make stop
```