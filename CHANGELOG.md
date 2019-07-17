## UNRELEASED

Refactor. EVM-Lite becomes a library

SECURITY:
- service: removed keystore
- service: removed /tx endpoint for unsigned txs
- service: removed CORS

FEATURES:

IMPROVEMENTS:

BUG FIXES:

## 0.2.0 (July 15, 2019)

SECURITY:

FEATURES:

IMPROVEMENTS:
- demo: Use evm-lite-lib package in demo scripts.
- state: Move genesis account creation from service to state. 
- state: PoA smart-contract bindings.

BUG FIXES:
- state: Initialize from empty state instead of latest trie root. This enables
         bootstrapping evm-lite/babble nodes from the babble DB only.

## V0.1.1 (January 28, 2019)

SECURITY:

FEATURES:

IMPROVEMENTS:
- deps: Use Geth v1.8.17
- consensus: Version 0.4.1 of Babble

BUG FIXES:
- deps: Update version of Geth. Version 1.7.0 had broken dependencies.

## V0.1.0 (October 14, 2018)

FEATURES:
- state: EVM, Trie, and Database.
- service: account management and HTTP API.
- consensus: simple consensus interface.
- consensus/solo: consensus implementation that simply relays transactions from
  service to state.
- consensus/babble: consensus implementation that uses an inmemory Babble node.
- consensus/raft: consensus implementation using hashicorp/raft
- engine: agent coordinating State, Service and Consensus.
- cmd: CLI
- deploy: scripts to create testnets locally or in AWS.
