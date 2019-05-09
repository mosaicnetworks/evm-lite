const util = require('util');
const fs = require('fs');
const JSONbig = require('json-bigint');
const argv = require('minimist')(process.argv.slice(2));
const evmlc = require('evm-lite-lib');
const solc = require('solc');
const prompt = require('prompt');

const FgRed = '\x1b[31m';
const FgGreen = '\x1b[32m';
const FgYellow = '\x1b[33m';
const FgBlue = '\x1b[34m';
const FgMagenta = '\x1b[35m';
const FgCyan = '\x1b[36m';
const FgWhite = '\x1b[37m';

const log = (color, text) => {
	console.log(color + text + '\x1b[0m');
};


const online = false;

const schema = {
    properties: {
      enters: {
        description: 'PRESS ENTER TO CONTINUE',
        ask: function() {
			 if (!online){console.log("Skipping prompt");}
          return online;
        }
      }
    }
  };


const step = message => {
	log(FgWhite, '\n' + message);
	return new Promise(resolve => {
		prompt.get(schema, function(err, res) {
			resolve();
		});
	});
};

const hardstep = message => {
	log(FgWhite, '\n' + message);
	return new Promise(resolve => {
		prompt.get('PRESS ENTER TO CONTINUE', function(err, res) {
			resolve();
		});
	});
};

const explain = message => {
	log(FgCyan, util.format('\nEXPLANATION:\n%s', message));
};

const space = () => {
	console.log('\n');
};

const sleep = function(time) {
	return new Promise(resolve => setTimeout(resolve, time));
};

/**
 * Demo starts here.
 */

function Node(name, host, port) {
	this.name = name;
	this.api = new evmlc.EVMLC(host, port, {
		from: '',
		gas: 1000000,
		gasPrice: 0
	});
	this.account = {};
}

const allAccounts = [];
const allNodes = [];

var poa_Event2 = {};
var contractPath = '';

const init = async () => {
	console.group('Initialize Nodes: ');

	const ips = argv.ips
		.replace(/\s/g, '')
		.split(',')
		.sort();
	const port = argv.port;
	const keystore = new evmlc.Keystore(argv.keystore, 'keystore');
	const passwordPath = argv.pwd;
	const password = readPasswordFile(passwordPath);

	contractPath = argv.contract;

	console.log('Sorted IPs: ', ips);
	console.log('Keystore Path: ', keystore.path);
	console.log('Password File Path: ', passwordPath);
	console.log('Contract File Path: ', contractPath);

	for (i = 0; i < ips.length; i++) {
		node = new Node(util.format('node%d', i + 1), ips[i], port);
		allNodes.push(node);
	}

	console.groupEnd();

	return {
		keystore,
		password
	};
};

const readPasswordFile = path => {
	return fs.readFileSync(path, { encoding: 'utf8' });
};

const decryptAccounts = async ({ keystore, password }) => {
	console.group('Decrypt Accounts');
	console.log('Password: ', password);

        const baseAccounts = await keystore.list(allNodes[0].api);

	for (const baseAccount of baseAccounts) {
		account = await keystore.decrypt(baseAccount.address, password);
		account.balance = baseAccount.balance;

		console.log('Decrypted: ', `${account.address}(${account.balance})`);
		allAccounts.push(account);
	}

	for (i = 0; i < allNodes.length; i++) {
		allNodes[i].api.defaultFrom = allAccounts[i].address;
		allNodes[i].account = allAccounts[i];
	}

	console.groupEnd();
};

const displayAllBalances = async () => {
	console.group('Current Account Balances');

	for (const node of allNodes) {
		baseAccount = await node.api.accounts.getAccount(node.account.address);
		console.log(`${node.name}: `, '\n', baseAccount, '\n');
	}
	console.groupEnd();
};

const transferRaw = async (from, to, value) => {
	console.group('Transfer Signed Locally');

	const transaction = await from.api.accounts.prepareTransfer(
		to.account.address,
		value
	);
	console.log('Transaction: ', transaction.parse(), '\n');

	await transaction.submit(from.account, { timeout: 3 });

	console.log('Receipt: ', await transaction.receipt);
};

const compiledSmartContract = async () => {
	const input = fs.readFileSync(contractPath, {
		encoding: 'utf8'
	});
	const output = solc.compile(input.toString(), 1);
	

   console.log("Output: ", output);
	const bytecode = output.contracts[`:POA_Event2`].bytecode;

   console.log("ByteCode: ", bytecode);

	const abi = output.contracts[`:POA_Event2`].interface;

	const contract = await allNodes[0].api.contracts.load(JSON.parse(abi), {
		data: bytecode
	});

	return contract;
};




class POA_Event2 {
	constructor(contract, account) {
		this.contract = contract;
		this.account = account;
	}

	async deploy(value) {
		await this.contract.deploy(this.account, [value], { timeout: 3 });
		console.log('Receipt:', this.contract.receipt);
		return this;
	}

	async genericTrans(functionname, value, params) {

		console.log(functionname, params);
		const transaction = await this.contract.methods[functionname](...params);

		transaction.value(value);

		console.log('Transaction: ', transaction.parse(), '\n');
		const response = await transaction.submit(this.account, { timeout: 3 });

      console.log('Response: ', response, '\n');
		if (value > 0)
		{
    		const receipt = await transaction.receipt;
   
   		const logs = this.contract.parseLogs(receipt.logs);
   
   		for (const log of logs) {
   			 console.log(
   					log.event || 'No Event Name',
   					JSON.stringify(log.args, null, 2)
   			 );
   		}
      }
		return transaction;
	}


	async genericTransAccount(fAccount, functionname, value, params) {

		console.log(functionname, params);
		const transaction = await this.contract.methods[functionname](...params);
      console.log("Created Transaction");
		transaction.value(value);

		console.log('Transaction: ', transaction.parse(), '\n');

		console.log('Next line fails:');

		const response = await transaction.submit(fAccount, { timeout: 3 });
		
      console.log('Transaction Submit Response: ', response, '\n');
		if (value > 0)
		{
    		const receipt = await transaction.receipt;
   
   		const logs = this.contract.parseLogs(receipt.logs);
   
   		for (const log of logs) {
   			 console.log(
   					log.event || 'No Event Name',
   					JSON.stringify(log.args, null, 2)
   			 );
   		}
      }
		return transaction;
	}




   async submitNominee(nomineeAddress, moniker, value) {

      console.log("submitNominee call");
      const transaction = await this.contract.methods.submitNominee(nomineeAddress, moniker);

      console.log("submitNominee call");

      transaction.value(value);

      console.log('Transaction: ', transaction.parse(), '\n');
      await transaction.submit(this.account, { timeout: 3 });

      const receipt = await transaction.receipt;
      const logs = this.contract.parseLogs(receipt.logs);

      for (const log of logs) {
         console.log(
            log.event || 'No Event Name',
            JSON.stringify(log.args, null, 2)
         );
      }


			return transaction;
		}



   async castNomineeVote(nomineeAddress, accepted, value) {

      const transaction = await this.contract.methods.castNomineeVote(nomineeAddress, accepted);
      transaction.value(value);

      console.log('Transaction: ', transaction.parse(), '\n');
      await transaction.submit(this.account, { timeout: 9 });

      const receipt = await transaction.receipt;
      const logs = this.contract.parseLogs(receipt.logs);

      for (const log of logs) {
         console.log(
            log.event || 'No Event Name',
            JSON.stringify(log.args, null, 2)
         );
      }


			return transaction;
		}

}

init()
	.then(object => decryptAccounts(object))
	.then(() => step('STEP 1) Get ETH Accounts'))
	.then(() => {
		space();
		return displayAllBalances();
	})
	.then(() =>
		explain(
			'Each node controls one account which allows it to send and receive Ether. \n' +
			'The private keys reside locally and directly on the evm-light nodes. In a \n' +
			'production setting, access to the nodes would be restricted to the people  \n' +
			'allowed to sign messages  private key. We also keep a local copy \n' +
			'of all the private keys to demonstrate client-side signing.'
		)
	)
//	.then(() => step('STEP 2) Send 500 wei (10^-18 ether) from node1 to node2'))
//	.then(() => {
//		space();
//		return transferRaw(allNodes[0], allNodes[1], 500);
//	})
//	.then(() =>
//		explain(
//			'We created an EVM transaction to send 500 wei from node1 to node2. The \n' +
//			'transaction was signed localy with node1 \'s private key and sent through node1. \n' +
//			'The client-facing service running in EVM-Lite relayed the transaction to Babble \n' +
//	 	  'for consensus ordering. Babble gossiped the raw transaction to the other Babble \n' +
//			'nodes which ran it through the consensus algorithm before committing it back to \n' +
//			'EVM-Lite as part of Block. So each node received and processed the transaction. \n' +
//			'They each applied the same changes to their local copy of the ledger.\n'
//		)
//	)
//	.then(() => step('STEP 3) Check balances again'))
//	.then(() => {
//		space();
//		return displayAllBalances();
//	})
//	.then(() =>
//		explain('Notice how the balances of node1 and node2 have changed.')
//	)
	.then(() =>
		step(
			'STEP 4) Deploy a POA_Event2 SmartContract for 1000 wei from node 1'
		)
	)
	.then(() => {
		space();
		return compiledSmartContract('PATH_TO_CONTRACT_HERE');
	})


	.then(async contract => {
		poa_Event2 = new POA_Event2(contract, allNodes[0].account);
		await poa_Event2.deploy(1000);
	})
	.then(() =>
		explain(
			'Here we compiled and deployed the POA_Event2 SmartContract. \n' +
			'The contract was written in the high-level Solidity language which compiles \n' +
			'down to EVM bytecode. To deploy the SmartContract we created an EVM transaction \n' +
			"with a 'data' field containing the bytecode. After going through consensus, the \n" +
			'transaction is applied on every node, so every participant will run a copy of \n' +
			'the same code with the same data.'
		)
	)


   .then(() => step('STEP 5.5) Check Number of Nodes in the White List'))
   .then(() => {
      space();
      return poa_Event2.genericTrans('getWhiteListCount', 0, []);
   })
   .then(() =>
      explain(
         "Get number of nodes in the whitelist. It should be one at this point"
      )
   )



	.then(() => step('STEP 8) Check if Node 2 is in nominee list'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('isNominee', 0, [allNodes[1].account.address]);
	})
	.then(() =>
		explain(
			'Should return a boolean - should be whitelisted here \n' +
			'.'
		)
	)


	.then(() => step('STEP 5) Nominate Node 2'))
	.then(() => {
		space();
	return poa_Event2.submitNominee(allNodes[1].account.address,'0x4e69636b00000000000000000000000000000000000000000000000000000000', 10000 );
	})
	.then(() =>
		explain(
			"Node 1 nominates node 2."
		)
	)



	.then(() => step('STEP 8) Check if Node 2 is in nominee list'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('isNominee', 0, [allNodes[1].account.address]);
	})
	.then(() =>
		explain(
			'Should return a boolean - should be whitelisted here \n' +
			'.'
		)
	)


	.then(() => step('STEP 6) Node 1 votes for Node 2'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('castNomineeVote', 20000, [allNodes[1].account.address, true]);
	})
	.then(() =>
		explain(
			'Node 1 casts a vote for node 2 \n' +
			'.'
		)
	)

	.then(() => step('STEP 7) Check if Node 2 is in whitelist'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('isWhitelisted', 0, [allNodes[1].account.address]);
	})
	.then(() =>
		explain(
			'Should return a boolean - should be whitelisted here \n' +
			'.'
		)
	)


	.then(() => step('STEP 8) Check if Node 2 is in nominee list'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('isNominee', 0, [allNodes[1].account.address]);
	})
	.then(() =>
		explain(
			'Should return a boolean - should be whitelisted here \n' +
			'.'
		)
	)


	.then(() => step('STEP 5.5) Check Number of Nodes in the White List'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('getWhiteListCount', 0, []);
	})
	.then(() =>
		explain(
			"Get number of nodes in the whitelist"
		)
	)



	.then(() => step('STEP 8) Check if Node 4 is in nominee list'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('isNominee', 0, [allNodes[3].account.address]);
	})
	.then(() =>
		explain(
			'Should return a boolean - should be whitelisted here \n' +
			'.'
		)
	)


	.then(() => step('STEP 5) Nominate Node 4'))
	.then(() => {
		space();
	return poa_Event2.submitNominee(allNodes[3].account.address,'0x4e69636b00000000000000000000000000000000000000000000000000000000', 10000 );
	})
	.then(() =>
		explain(
			"Node 1 nominates node 4."
		)
	)



	.then(() => step('STEP 8) Check if Node 4 is in nominee list'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('isNominee', 0, [allNodes[3].account.address]);
	})
	.then(() =>
		explain(
			'Should return a boolean - should be whitelisted here \n' +
			'.'
		)
	)


	.then(() => step('STEP 6) Node 2 votes for Node 4'))
	.then(() => {
		space();
		return poa_Event2.genericTransAccount(allNodes[1].account, 'castNomineeVote', 20000, [allNodes[3].account.address, true]);
	})
	.then(() =>
		explain(
			'Node 2 casts a vote for node 4 \n' +
			'.'
		)
	)

	.then(() => step('STEP 7) Check if Node 4 is in whitelist'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('isWhitelisted', 0, [allNodes[3].account.address]);
	})
	.then(() =>
		explain(
			'Should return a boolean - should be whitelisted here \n' +
			'.'
		)
	)


	.then(() => step('STEP 8) Check if Node 4 is in nominee list'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('isNominee', 0, [allNodes[3].account.address]);
	})
	.then(() =>
		explain(
			'Should return a boolean - should be whitelisted here \n' +
			'.'
		)
	)


	.then(() => step('STEP 5.5) Check Number of Nodes in the White List'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('getWhiteListCount', 0, []);
	})
	.then(() =>
		explain(
			"Get number of nodes in the whitelist"
		)
	)



	.then(() => step('STEP 6) Node 1 votes for Node 4'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('castNomineeVote', 20000, [allNodes[3].account.address, true]);
	})
	.then(() =>
		explain(
			'Node 1 casts a vote for node 4 \n' +
			'.'
		)
	)

	.then(() => step('STEP 7) Check if Node 4 is in whitelist'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('isWhitelisted', 0, [allNodes[3].account.address]);
	})
	.then(() =>
		explain(
			'Should return a boolean - should be whitelisted here \n' +
			'.'
		)
	)


	.then(() => step('STEP 8) Check if Node 4 is in nominee list'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('isNominee', 0, [allNodes[3].account.address]);
	})
	.then(() =>
		explain(
			'Should return a boolean - should be whitelisted here \n' +
			'.'
		)
	)


	.then(() => step('STEP 5.5) Check Number of Nodes in the White List'))
	.then(() => {
		space();
		return poa_Event2.genericTrans('getWhiteListCount', 0, []);
	})
	.then(() =>
		explain(
			"Get number of nodes in the whitelist"
		)
	)


	.catch(err => log(FgRed, err));
