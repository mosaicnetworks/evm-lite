## UNRELEASED

SECURITY:

FEATURES:

IMPROVEMENTS:

BUG FIXES:

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
