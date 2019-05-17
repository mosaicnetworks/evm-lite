# Authority and Babble Tests


##Set Up
We can pre-load a configuration with the smart contract embedded in the genesis block.

```bash
$ cd [...evm-lite]/deploy
$ make prebuilt PREBUILT=babblepoa
```

We can then start this test network with:

```bash
$ make start CONSENSUS=babble NODES=4
# We cheat and copy the IP config - the later docker commands wipe it.
$ cp terraform/local/ips.dat ../demo
```

Kill nodes 2 and 3. This will not be necessary longer term as the evm-lite instance will abort if it lacks the authority to join the network
```bash
$ docker exec -it node2 kill 1
$ docker exec -it node3 kill 1
```



## Demo Workflow

|Step|General|Node 0 | Node 1 | Node 2|Node3|
|  --:|------|------|------|------|------|
|1|Start Network as per above | | | || 
|2| |Nominate Node 2 | | | | 
|3| |Vote for Node 2 | | | | 
|4| | | |Try to Join. Should Fail | |
|5| | |Vote for Node  2. Should announce result | | | 
|6| | | |Try to Join. Should Succeed | | 
|7| |Nominate Node 3 | | | |  
|8| |Vote for Node 3 | | | |  
|9| | |Vote for Node 3 | | |  
|10| | | |Vote Against Node 3, Should reach decision| |  
|11| | | ||Try to Join. Should fail |  



## Enhancements
Sketching out this demo has made it apparent that nominating someone then having to vote for them is possibly a nonsense.

## Some Docker Commands


To stop and restart a running node:
```bash
$ docker exec -it node0 ps 	-fe
PID   USER     TIME  COMMAND
    1 1000      0:08 evml run babble
   45 1000      0:00 ps -fe
$ docker exec -it node0 kill 1
$ docker ps --all
CONTAINER ID        IMAGE                           COMMAND             CREATED             STATUS                      PORTS                NAMES
6a28ef4d4b24        mosaicnetworks/evm-lite:0.2.0   "evml run babble"   15 hours ago        Up 15 hours                 8000/tcp, 8080/tcp   node2
a4aea452898e        mosaicnetworks/evm-lite:0.2.0   "evml run babble"   15 hours ago        Up 15 hours                 8000/tcp, 8080/tcp   node1
cf9323997329        mosaicnetworks/evm-lite:0.2.0   "evml run babble"   15 hours ago        Up 15 hours                 8000/tcp, 8080/tcp   node3
9c4291085800        mosaicnetworks/evm-lite:0.2.0   "evml run babble"   15 hours ago        Exited (2) 22 seconds ago                        node0
$ docker restart node0
node0
$ docker exec -it node0 ps -fe
PID   USER     TIME  COMMAND
    1 1000      0:01 evml run babble
   24 1000      0:00 ps -fe

```
