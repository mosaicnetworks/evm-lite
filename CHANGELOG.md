
# Changelog

## Unreleased

SECURITY:
FEATURES:
IMPROVEMENTS:
BUG FIXES:

## v0.3.7 (November 27, 2019)

BUG FIXES:

- state: set account storage and nonce when creating accounts

## v0.3.6 (November 7, 2019)

IMPROVEMENTS:

- service: new export function and endpoint that returns a full snapshot of the
           state, which can be reused as a genesis file. 

## v0.3.5 (October 15, 2019)

IMPROVEMENTS:

- state: optimisations and performance tuning

## v0.3.4 (September 18, 2019)

IMPROVEMENTS:

- state: more granular use of mutexes.
- service: higher throughput thanks to above improvement.

## v0.3.3 (September 13, 2019)

FEATURES:

- currency: new denominations for token units

BUG FIXES:

- state: handling transaction promises and errors

## v0.3.2 (September 6, 2019)

FEATURES:
- state: make use of coinbase address
- service: min gas price

## v0.3.1 (August 29, 2019)

SECURITY:
- service: enable CORS

IMPROVEMENTS:
- service: make /tx synchronous (directly return receipt)

## v0.3.0 (August 8, 2019)

Refactor. EVM-Lite becomes a library

SECURITY:
- service: removed keystore
- service: removed /tx endpoint for unsigned txs
- service: removed CORS

FEATURES:
- service: handlers registered to http.DefaultServerMux

## v0.2.1 (July 22, 2019)

BUG FIXES: 
- service: Force O gasPrice on readonly transactions (issue 5 of monetd)

## v0.2.0 (July 22, 2019)

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
