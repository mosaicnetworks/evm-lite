#Demo

The scripts described below attempt to run the entire demo process via bash and evmlc.

##Installation

First we fire up our prebuilt instances.

```bash
$ cd [evm-lite]/deploy
$ make prebuilt PREBUILT=babblepoademo
$ make start CONSENSUS=babble NODES=10
```

Now we shutdown nodes 4 to 9 as they are not authorised at this time. We can restart them when we need them. 

We can use docker commands directly:

```bash
$ for i in  {4..9}; do docker exec node${i} kill 1; done
```

Or we can use a script:

```bash
$ scripts/kill_docker_node.sh node4 node5 node6 node7 node8 node9
```


##Firing up a Node

For this example we will be firing up node6. The genesis peers files are all preset to the initial 4 nodes and thus do not need to be amended.

```bash
$ cd [evm-lite]/deploy
$ scripts/set_docker_peers.sh node6
$ scripts/restart_docker_node.sh node6
```

The above scripts get a list of running nodes via a docker ps command. It picks a random node from the list, and queries the peers endpoint to get a current list of peers. That list is then written to the docker container. The node is then restarted and picked up the revised node configuration. 

## Our Script Toolkit

We do not amend the path, but they is a set of scripts to ease node manipulation in deploy/scripts.
```bash
$ cd [evm-lite]/deploy/scripts
```

This scripts calls all active nodes to vote for the given node. 
```bash
$ all_vote_node.sh node5
```



This script detects the current docker nodes and picks one at random to request a peers list from. The script takes a single argument of the output file to write the peers to. 
```bash
$ get_current_peers.sh /tmp/peers.json
```

Script to return the IP address of a container. It is derived from the docker container, not the ips.dat file.
```bash
$ get_ip_of_container.sh node0
```

Returns the ethereum address associated with the validator of a node. 
```bash
$ get_node_address.sh node4
```

Returns a list of running nodes and their IP addresses - using the mapping from ips.dat
```bash
$ get_running_nodes.sh
```

Gets the ethereum address for the validator for a sleeping node. It pulls the data from key.json in the docker container.
```bash
$ get_sleeping_node_address.sh node0
```

Kills a docker node by killing the evml process. 
```bash
$ kill_docker_node.sh node9
```

Makes node0 nominate node9. Wrapper for evmlc. 
```bash
$ nominate_node.sh node0 node9
```

This script retrieves a list of nodes currently active, gets their IP address, and rewrites to ips.dat. This is intended to address issues with nodes changing their ip as they are stopped and restarted. 
```bash
$ rebuild_ipsdat.sh
```

Start a docker node.
```bash
$ restart_docker_node.sh node8
```

Gets the current peers file from a live node and copies it to however many nodes are specified on the command line. This is intended to allow an exited node be started with a new peer set. The genesis peers are all set in the initial node creation. 
```bash
$ set_docker_peers.sh node5 node6
```

This script makes node0 vote for node3. Change the 3rd parameter to false to vote against. 
```bash
$ vote_node.sh node0 node3 true
```


This script writes a new configuration file to ~/.evmlc/config.toml
```bash
$ write_evmlc_config.sh node1
```


# The Demo

Run the commands in the installation section above. This will leave us with Nodes 0 to 3 all running. Those 4 nodes associated with the accounts that are baked into the POA Solidity contract. 

We can peek into a container to get that list with the following command:
```bash
$ docker exec node0 grep address /home/1000/.evm-lite/eth/genesis.json
            "address": "abbaabbaabbaabbaabbaabbaabbaabbaabbaabba",
                  "address": "0x815A9d3C1b9b2Ec8f49F5730830004BD2F83b8b8",
                  "address": "0x8c266894Ac9f23e4cF5300220dfe79896E7576fE",
                  "address": "0x910e467DF064083407ceb8406840a20c2f25DCc2",
                  "address": "0x01728D4D07838E7A5DF5d45435272a9592C5ea4d",

```
The first address is the contract address, abba recurring. The following items are baked into the contract. The address whose credentials are stored within a running node can be displayed with the command below:

```bash
$ ./scripts/get_node_address.sh node0
"address":"815a9d3c1b9b2ec8f49f5730830004bd2f83b8b8"
```
You can verify that each of the nodes 0 to 3 has the credentials for a node on the initial white list. 

##Adding a Node

To add a node the following steps are undertaken:

+ The node is proposed by an existing validator
+ Each validator votes for them
+ The node joins

```bash
$ cd [evm-lite]/deploy
# Set to node 0
$ scripts/write_evmlc_config.sh node0
# Nominate
$ scripts/nominate_node.sh node0 node4
$ scripts/vote_node.sh node0 node4 true

# node1 votes for node4
$ scripts/write_evmlc_config.sh node1
$ scripts/vote_node.sh node1 node4 true

# node2 votes for node4
$ scripts/write_evmlc_config.sh node2
$ scripts/vote_node.sh node2 node4 true

# node3 votes for node4
$ scripts/write_evmlc_config.sh node2
$ scripts/vote_node.sh node2 node4 true

# It is now unanimous and node 4 is added to the whitelist

# First we set the peers to the current values for node 4
$ scripts/set_docker_peers.sh node4
# Then we restart the node
$ scripts/restart_docker_node.sh node4
# Then we regenerate ips.dat
$ scripts/rebuild_ipsdat.sh
```



There is a briefer version using all_vote_node.sh to wrap the individual voting calls up into a single script.
```bash
$ cd [evm-lite]/deploy
# Set to node 0
$ scripts/write_evmlc_config.sh node0
# Nominate
$ scripts/nominate_node.sh node0 node4
$ scripts/all_vote_node.sh node4

# It is now unanimous and node 4 is added to the whitelist

# First we set the peers to the current values for node 4
$ scripts/set_docker_peers.sh node4
# Then we restart the node
$ scripts/restart_docker_node.sh node4
# Then we regenerate ips.dat
$ scripts/rebuild_ipsdat.sh
```

## Removing a Node

This is not currently within scope.  Whilst there is a leave event defined in Babble, we currently do not expose a leave method through evmlc. 




