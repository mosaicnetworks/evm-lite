## UNRELEASED

FEATURES:
- state: EVM, Trie, and Database.
- service: account management and HTTP API.
- consensus: simple consensus interface.
- consensus/solo: consensus implmentation that simply relays transactions from
  service to state.
- consensus/babble: consensus implementation that uses an inmemory Babble node.
- engine: agent that coordinates State, Service and Consensus.
- cmd: CLI
- deploy: scripts to deploy local docker network