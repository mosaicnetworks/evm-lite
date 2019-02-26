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

const step = message => {
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

//------------------------------------------------------------------------------

const sleep = function(time) {
	return new Promise(resolve => setTimeout(resolve, time));
};

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
var crowdFunding = {};

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

	console.log('Sorted IPs: ', ips);
	console.log('Keystore Path: ', keystore.path);
	console.log('Password File Path: ', passwordPath);

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

	const baseAccounts = await keystore.list(true, allNodes[0].api);

	for (const baseAccount of baseAccounts) {
		account = await keystore.decryptAccount(baseAccount.address, password);
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
		console.log(node.name, JSON.stringify(baseAccount, null, 2));
	}
	console.groupEnd();
};

const transferRaw = async (from, to, value) => {
	console.group('Transfer Signed Locally');

	const transaction = await from.api.accounts.prepareTransfer(
		to.account.address,
		value
	);
	console.log('Transaction: ', transaction.parse());

	await transaction.submit({ timeout: 2 }, from.account);

	console.log('Receipt: ', await transaction.receipt);
};

const compiledSmartContract = async () => {
	const input = fs.readFileSync('../smart-contracts/CrowdFunding.sol', {
		encoding: 'utf8'
	});
	const output = solc.compile(input.toString(), 1);
	const bytecode = output.contracts[`:CrowdFunding`].bytecode;
	const abi = output.contracts[`:CrowdFunding`].interface;

	const contract = await allNodes[0].api.contracts.load(JSON.parse(abi), {
		data: bytecode
	});

	return contract;
};

class CrowdFunding {
	constructor(contract, account) {
		this.contract = contract;
		this.account = account;
	}

	async deploy(value) {
		await this.contract.deploy(this.account, [value], { timeout: 2 });
		console.log('Receipt:', this.contract.receipt);
		return this;
	}

	async contribute(value) {
		const transaction = await this.contract.methods.contribute();
		transaction.value(value);

		console.log('Transaction: ', transaction.parse());
		await transaction.submit({ timeout: 2 }, this.account);

		console.log('Receipt: ', await transaction.receipt);
		return transaction;
	}

	async checkGoalReached() {
		const transaction = await this.contract.methods.checkGoalReached();
		const response = await transaction.submit({ timeout: 2 }, this.account);

		console.log('Response: ', response);
		return response;
	}

	async settle() {
		const transaction = await this.contract.methods.settle();

		console.log('Transaction: ', transaction.parse());
		await transaction.submit({ timeout: 2 }, this.account);

		console.log('Receipt: ', await transaction.receipt);
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
				'The private keys reside directly on the evm-babble nodes. In a production \n' +
				'setting, access to the nodes would be restricted to the people allowed to \n' +
				'sign messages with the private key. We also keep a local copy of all the private \n' +
				'keys to demonstrate client-side signing.'
		)
	)
	.then(() => step('STEP 2) Send 500 wei (10^-18 ether) from node1 to node2'))
	.then(() => {
		space();
		return transferRaw(allNodes[0], allNodes[1], 500);
	})
	.then(() =>
		explain(
			'We created an EVM transaction to send 500 wei from node1 to node2. The \n' +
				'transaction was sent to node1 which controls the private key for the sender. \n' +
				'EVM-Babble converted the transaction into raw bytes, signed it and submitted \n' +
				'it to Babble for consensus ordering. Babble gossiped the raw transaction to \n' +
				'the other Babble nodes which ran it through the consensus algorithm until they \n' +
				'were each ready to commit it back to EVM-BABBLE. So each node received and \n' +
				'processed the transaction. They each applied the same changes to their local \n' +
				'copy of the ledger.'
		)
	)
	.then(() => step('STEP 3) Check balances again'))
	.then(() => {
		space();
		return displayAllBalances();
	})
	.then(() =>
		explain('Notice how the balances of node1 and node2 have changed.')
	)
	.then(() =>
		step(
			'STEP 6) Deploy a CrowdFunding SmartContract for 1000 wei from node 1'
		)
	)
	.then(() => {
		space();
		return compiledSmartContract();
	})
	.then(async contract => {
		crowdFunding = new CrowdFunding(contract, allNodes[0].account);
		await crowdFunding.deploy(1000);
	})
	.then(() =>
		explain(
			'Here we compiled and deployed the CrowdFunding SmartContract. \n' +
				'The contract was written in the high-level Solidity language which compiles \n' +
				'down to EVM bytecode. To deploy the SmartContract we created an EVM transaction \n' +
				"with a 'data' field containing the bytecode. After going through consensus, the \n" +
				'transaction is applied on every node, so every participant will run a copy of \n' +
				'the same code with the same data.'
		)
	)
	.then(() => step('STEP 4) Contribute 499 wei from node 1'))
	.then(() => {
		space();
		return crowdFunding.contribute(499);
	})
	.then(() =>
		explain(
			"We created an EVM transaction to call the 'contribute' method of the SmartContract. \n" +
				"The 'value' field of the transaction is the amount that the caller is actually \n" +
				'going to contribute. The operation would fail if the account did not have enough Ether. \n' +
				'As an exercise you can check that the transaction was run through every Babble \n' +
				"node and that node2's balance has changed."
		)
	)
	.then(() => step('STEP 5) Check goal reached'))
	.then(() => {
		space();
		return crowdFunding.checkGoalReached();
	})
	.then(() =>
		explain(
			'Here we called another method of the SmartContract to check if the funding goal \n' +
				'was met. Since only 499 of 1000 were received, the answer is no.'
		)
	)
	.then(() => step('STEP 6) Contribute 501 wei from node 1 again'))
	.then(() => {
		space();
		return crowdFunding.contribute(501);
	})
	.then(() => step('STEP 7) Check goal reached'))
	.then(() => {
		space();
		return crowdFunding.checkGoalReached();
	})
	.then(() =>
		step(
			'STEP 8) Before we `settle` lets check balances again to show that node 1 balance decreased by a total of 1000.'
		)
	)
	.then(() => {
		space();
		return displayAllBalances();
	})
	.then(() =>
		explain('Since the funding goal was reached we can now settle.')
	)
	.then(() => step('STEP 9) Settle'))
	.then(() => {
		space();
		return crowdFunding.settle();
	})
	.then(() =>
		explain(
			'The funds were transferred from the SmartContract back to node1.'
		)
	)
	.then(() => step('STEP 10) Check balances again'))
	.then(() => {
		space();
		return displayAllBalances();
	})
	.catch(err => log(FgRed, err));
