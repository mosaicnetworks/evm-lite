# EVM-LITE

[![CircleCI](https://circleci.com/gh/mosaicnetworks/evm-lite.svg?style=svg)](https://circleci.com/gh/mosaicnetworks/evm-lite)
[![Documentation Status](https://readthedocs.org/projects/evm-lite/badge/?version=latest)](https://evm-lite.readthedocs.io/en/latest/?badge=latest)
[![Go Report](https://goreportcard.com/badge/github.com/mosaicnetworks/evm-lite)](https://goreportcard.com/report/github.com/mosaicnetworks/evm-lite)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## A lean Ethereum node with interchangeable consensus.

We took the [Go-Ethereum](https://github.com/ethereum/go-ethereum)
implementation (Geth) and extracted the EVM and Trie components to create a
lean and modular version with interchangeable consensus.

The EVM is a virtual machine specifically designed to run untrusted code on a
network of computers. Every transaction applied to the EVM modifies the State
which is persisted in a Merkle Patricia tree. This data structure allows to
simply check if a given transaction was actually applied to the VM and can
reduce the entire State to a single hash (merkle root) rather analogous to a
fingerprint.

The EVM is meant to be used in conjunction with a system that broadcasts
transactions across network participants and ensures that everyone executes the
same transactions in the same order. Ethereum uses a Blockchain and a Proof of
Work consensus algorithm. EVM-Lite makes it easy to use any consensus system,
including [Babble](https://github.com/mosaicnetworks/babble) .

## ARCHITECTURE

```
                +-------------------------------------------+
+----------+    |  +-------------+         +-------------+  |       
|          |    |  | Service     |         | State       |  |
|  Client  <-----> |             | <------ |             |  |
|          |    |  | -API        |         | -EVM        |  |
+----------+    |  |             |         | -Trie       |  |
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

## Usage

EVM-Lite is a Go library, which is meant to be used in conjunction with a 
consensus system like Babble, Tendermint, Raft, etc.

This repo contains **Solo**, a bare-bones implementation of the consensus 
interface, which is used for testing or launching a standalone node. It relays
transactions directly from Service to State.

## Configuration

The Ethereum genesis file defines Ethereum accounts and is stripped of all the 
Ethereum POW stuff. This file is useful to predefine a set of accounts that own 
all the initial Ether at the inception of the network.

Example Ethereum genesis.json defining two account:
```json
{
   "alloc": {
        "629007eb99ff5c3539ada8a5800847eacfc25727": {
            "balance": "1337000000000000000000"
        },
        "e32e14de8b81d8d3aedacb1868619c74a68feab0": {
            "balance": "1337000000000000000000"
        }
   }
}
```

## Database

EVM-Lite will use a LevelDB database to persist state objects. The file of the  
database can be specified with the `db` flag which defaults to
`<datadir>/eth/chaindata`.  

## API
The Service exposes an API at the address specified by the [XXX address config]
for clients to interact with Ethereum.  

### Get account

This method retrieves the information about any account.  

```bash
host:~$ curl http://[api_addr]/account/0x629007eb99ff5c3539ada8a5800847eacfc25727 -s | json_pp
{
    "address":"0x629007eb99ff5c3539ada8a5800847eacfc25727",
    "balance":1337000000000000000000,
    "nonce":0
}
```

### Send raw signed transactions

This endpoint allows sending NON-READONLY transactions ALREADY SIGNED. The
client is left to compose a transaction, sign it and RLP encode it. The
resulting bytes, represented as a Hex string, are to this method to be forwarded 
to the EVM.

This is an ASYNCHRONOUS operation and the effect on the State should be verified
by fetching the transaction' receipt.

example:
```bash
host:~$ curl -X POST http://[api_addr]/rawtx -d '0xf8628080830f424094564686380e267d1572ee409368e1d42081562a8e8201f48026a022b4f68bfbd4f4c309524ebdbf4bac858e0ad65fd06108c934b45a6da88b92f7a046433c388997fd7b02eb7128f4d2401ef2d10d574c42edf15875a43ee51a1993' -s | json_pp
{
    "txHash":"0x5496489c606d74ad7435568393fa2c4619e64497267f80864109277631aa849d"
}
```

### Get Transaction receipt

example:
```bash
host:~$ curl http://[api_addr]/tx/0xeeeed34877502baa305442e3a72df094cfbb0b928a7c53447745ff35d50020bf -s | json_pp
{
   "to" : "0xe32e14de8b81d8d3aedacb1868619c74a68feab0",
   "root" : "0xc8f90911c9280651a0cd84116826d31773e902e48cb9a15b7bb1e7a6abc850c5",
   "gasUsed" : "0x5208",
   "from" : "0x629007eb99ff5c3539ada8a5800847eacfc25727",
   "transactionHash" : "0xeeeed34877502baa305442e3a72df094cfbb0b928a7c53447745ff35d50020bf",
   "logs" : [],
   "cumulativeGasUsed" : "0x5208",
   "contractAddress" : null,
   "logsBloom" : "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
}

```

## Get consensus info

The `/info` endpoint exposes a map of information provided by the consensus
system.

example (with Babble consensus):
```bash
host:-$ curl http://[api_addr]/info | json_pp
{
   "rounds_per_second" : "0.00",
   "type" : "babble",
   "consensus_transactions" : "10",
   "num_peers" : "4",
   "consensus_events" : "10",
   "sync_rate" : "1.00",
   "transaction_pool" : "0",
   "state" : "Babbling",
   "events_per_second" : "0.00",
   "undetermined_events" : "22",
   "id" : "1785923847",
   "last_consensus_round" : "1",
   "last_block_index" : "0",
   "round_events" : "0"
}

```

## CLIENT

Please refer to [EVM-Lite CLI](https://github.com/mosaicnetworks/evm-lite-cli)
for Javascript utilities and a CLI to interact with the API.

## DEV

DEPENDENCIES

We use glide to manage dependencies:

```bash
[...]/evm-lite$ curl https://glide.sh/get | sh
[...]/evm-lite$ glide install
```
This will download all dependencies and put them in the **vendor** folder; it
could take a few minutes.
