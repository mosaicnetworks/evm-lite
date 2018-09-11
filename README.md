# EVM-LITE
A simple wrapper for the Ethereum Virtual Machine and State Trie that plugs into 
different consensus systems.

## ARCHITECTURE

```
                +-------------------------------------------+
+----------+    |  +-------------+         +-------------+  |       
|          |    |  | Service     |         | State App   |  |
|  Client  <-----> |             | <------ |             |  |
|          |    |  | -API        |         | -EVM        |  |
+----------+    |  | -Keystore   |         | -Trie       |  |
                |  |             |         | -Database   |  |
                |  +-------------+         +-------------+  |
                |         |                       |         |
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

Consensus Implementations:

- **SOLO**: No Consensus. Transactions are relayed directly from Service to 
            State
- **BABBLE**: Inmemory Babble node.

## USAGE

Each consensus has its own subcommand `evm-lite [consensus]`

```
Usage:
  evm-lite [command]

Available Commands:
  babble      Run the evm-lite node with Babble consensus
  help        Help about any command
  solo        Run the evm-lite node with Solo consensus (no consensus)
  version     Show version info

Flags:
      --datadir string        Top-level directory for configuration and data (default "/home/martin/.evm-lite")
      --eth.api_addr string   Address of HTTP API service (default ":8080")
      --eth.cache int         Megabytes of memory allocated to internal caching (min 16MB / database forced) (default 128)
      --eth.db string         Eth database file (default "/home/martin/.evm-lite/eth/chaindata")
      --eth.genesis string    Location of genesis file (default "/home/martin/.evm-lite/eth/genesis.json")
      --eth.keystore string   Location of Ethereum account keys (default "/home/martin/.evm-lite/eth/keystore")
      --eth.pwd string        Password file to unlock accounts (default "/home/martin/.evm-lite/eth/pwd.txt")
  -h, --help                  help for evm-lite
      --log_level string      debug, info, warn, error, fatal, panic (default "debug")

Use "evm-lite [command] --help" for more information about a command.

```

## DEV

To add a new consensus system:

- implement the consensus interface (consensus/consensus.go)
- create the corresponding CLI subcommand in cmd/commands/
- register that command to the root command

