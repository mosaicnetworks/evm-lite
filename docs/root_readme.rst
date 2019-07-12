.. role:: raw-html-m2r(raw)
   :format: html


EVM-LITE
========

A lean Ethereum node with interchangeable consensus.
----------------------------------------------------

We took the `Go-Ethereum <https://github.com/ethereum/go-ethereum>`_
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
including `Babble <https://github.com/mosaicnetworks/babble>`_ .

ARCHITECTURE
------------

::

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



Consensus Implementations:
--------------------------


* 
  SOLO\ : No Consensus. Transactions are relayed directly from Service to
  State.

* 
  `BABBLE <https://github.com/mosaicnetworks/babble>`_\ : Inmemory Babble node.
  EVM-Lite does not support Babble's FastSync and Dynamic Membership protocols
  yet, so it is important to set the ``--store`` flag, and a high ``--sync-limit`` 
  value. 

* 
  `RAFT <https://github.com/hashicorp/raft>`_\ : Hashicorp implementation of
  Raft (limited).

more to come...
