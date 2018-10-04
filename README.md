# EVM-LITE
## Ethereum with interchangeable consensus.

We took the [Go-Ethereum](https://github.com/ethereum/go-ethereum) 
implementation (Geth) and extracted the EVM and Trie components to create a 
modular version with interchangeable consensus. 

## ARCHITECTURE

```
                +-------------------------------------------+
+----------+    |  +-------------+         +-------------+  |       
|          |    |  | Service     |         | State       |  |
|  Client  <-----> |             | <------ |             |  |
|          |    |  | -API        |         | -EVM        |  |
+----------+    |  | -Keystore   |         | -Trie       |  |
                |  |             |         | -Database   |  |
                |  +-------------+         +-------------+  |
                |         |                       ^         |     
                |         v                       |         |
                |  +-------------------------------------+  |
                |  | Engine                              |  |
                |  |                                     |  | 
                |  |       +----------------------+      |  |
                |  |       | Consensus            |      |  |
                |  |       +----------------------+      |  |
                |  |                                     |  |
                |  +-------------------------------------+  |
                |                                           | 
                +-------------------------------------------+

```

## Consensus Implementations:

- **SOLO**: No Consensus. Transactions are relayed directly from Service to 
            State

- **BABBLE**: Inmemory Babble node.

- **RAFT**: Hashicorp implementation of Raft (limited).

more to come...

## USAGE

Each consensus has its own subcommand `evml [consensus]`

```
Ethereum with interchangeable consensus

Usage:
  evml [command]

Available Commands:
  babble      Run the evm-lite node with Babble consensus
  help        Help about any command
  raft        Run the evm-lite node with Raft consensus
  solo        Run the evm-lite node with Solo consensus (no consensus)
  version     Show version info

Flags:
      --datadir string        Top-level directory for configuration and data (default "/home/martin/.evm-lite")
      --eth.cache int         Megabytes of memory allocated to internal caching (min 16MB / database forced) (default 128)
      --eth.db string         Eth database file (default "/home/martin/.evm-lite/eth/chaindata")
      --eth.genesis string    Location of genesis file (default "/home/martin/.evm-lite/eth/genesis.json")
      --eth.keystore string   Location of Ethereum account keys (default "/home/martin/.evm-lite/eth/keystore")
      --eth.listen string     Address of HTTP API service (default ":8080")
      --eth.pwd string        Password file to unlock accounts (default "/home/martin/.evm-lite/eth/pwd.txt")
  -h, --help                  help for evml
      --log string            debug, info, warn, error, fatal, panic (default "debug")

Use "evml [command] --help" for more information about a command.

```

Options can also be specified in a `evml.toml` file in the `datadir`. 

ex (evml.toml):
``` toml
log=info
[eth]
db = "/eth.db"
[babble]
listen="127.0.0.1:1337"
```

## DEV

DEPENDENCIES

We use glide to manage dependencies: 

```bash
[...]/evm-lite$ curl https://glide.sh/get | sh
[...]/evm-lite$ glide install
```
This will download all dependencies and put them in the **vendor** folder; it 
could take a few minutes.

CONSENSUS

To add a new consensus system:

- implement the consensus interface (consensus/consensus.go)
- add a property to the the global configuration object (config/config.go)
- create the corresponding CLI subcommand in cmd/evml/commands/
- register that command to the root command


## DEPLOY

We provide a set of scripts to automate the deployment of testnets. This 
requires [terraform](https://www.terraform.io/) and 
[docker](https://www.docker.com/).

Support for AWS is also available (cf. deploy/)



