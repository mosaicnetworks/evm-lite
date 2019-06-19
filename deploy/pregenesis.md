# Embedding Smart Contracts in the Genesis Block

There is a new compile process for creating the genesis.json file, which defines 
the Genesis Block. It can be invoked as follows:

It should be noted that the constructor never fires as the contract is not 
placed in the block via the usual mechanism. 

## pregenesis.json

The genesis block in a POA system is defined by the file **pregenesis.json**. In format it is identical to genesis.json with an additional **precompiler** node at the top level.


			{
				"precompiler":{
					"contracts": [
						{ 
							"address": "0xabababababababababababababababababababab",
							"filename": "genesis.sol",
							"contractname": "POA_Genesis",
							"authorising": true,
							"balance": "1337000000000000000001",
							"preauthorised": [
								{ "address": "abababababababababababababababababababab"},
								{ "address": "eb9d65b3358c1db0df3f5bafbffdd7a8d666d547"}								
							]
						}
					]				
				
				},
				"alloc": {
				"eb9d65b3358c1db0df3f5bafbffdd7a8d666d547": {
						"balance": "1337000000000000000000" 
					},
			      	"70ab04a5ee08a3db4481676b7a4b76e060a55ff8": {
						"balance": "1337000000000000000000" 
					}
				}
			}


Currently, the precompiler node contains a single contracts node that contains an array of contract objects, which contain the following keys:

- **address** - the ethereum address where the contract will be stored
- **filename** - the solidity source file of this contract. Currently with no path, it is assumed to be in the same directory as the pregenesis.json file
- **contractname** - the name of the Contract in Solidity - a solidity file may contain multiple contracts
- **authorising** - set to true if this contract handles authorisation for Babble-POA, false if it is just another contract.
- **balance** - the balance of this ethereum account. Not requried to place the contract, but convenient to set here if required.
- **preauthorised** - an array of validator objects, each of which that the following properties:
    - **address** - the ethereum address of the validator
    - **moniker** - a moniker for the validator 
    
## Dependencies
    
The process requires that solc be installed and available in your path. At this stage the command line version has been used rather than the node wrapper as it the wrapper does appear to expose all of the required options. For ubuntu (and many other Linux) it can be installed by:

	$ sudo apt-get install solc


NPM is also required. You can check requirements as below:

```bash
$ cd [evm-lite]/deploy
$ make checktools
```
    

## Build Process


```bash
$ cd [evm-lite]/deploy
$ make compile CONSENSUS=babble
```
      
### Output Files

The output directory (/deploy/conf/babble/conf) contains the following files:

+ **contract[*n*].sol** - the solidity code for the *nth* (zero based) contract. N.B. this version has the genesis whitelist values embedded within it. 
+ **contract[*n*].out** - details of the environment used to compile the *nth* contract. Currently contains a JSON object with a single key set to the compiler version. 
+ **POA_Genesis.abi** - the ABI output for the contract. POA_Genesis is the solidity contract name in this case. 
+ **POA_Genesis.bin-runtime** - the compiled byte code of the solidity contract. This may be removed in future as the same information is available in the genesis.json file.
+ **genesis.json** - the final genesis.json file

