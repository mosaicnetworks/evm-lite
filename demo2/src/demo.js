const util = require('util');
const fs = require('fs');
const JSONbig = require('json-bigint');
const argv = require('minimist')(process.argv.slice(2));
const evmlc = require('evm-lite-lib');
const solc = require('solc');

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

	await transaction.submit({}, from.account);

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

	async deploy() {
		await this.contract.deploy(this.account, [1000]);

		return this;
	}

	async contribute(value) {
		const transaction = await this.contract.methods.contribute();
		transaction.value(value);

		await transaction.submit({}, this.account);

		return transaction;
	}

	async checkGoalReached() {
		const transaction = await this.contract.methods.checkGoalReached();
		const response = await transaction.submit({}, this.account);

		return response;
	}
}

init()
	.then(object => decryptAccounts(object))
	.then(() => displayAllBalances())
	.then(() => transferRaw(allNodes[0], allNodes[1], 200))
	.then(() => displayAllBalances())
	.then(() => compiledSmartContract())
	.then(contract => new CrowdFunding(contract, allNodes[0].account))
	.then(contract => contract.deploy())
	.then(async contract => {
		const transaction = await contract.contribute(20);
		console.log(
			'Contribute Transaction Receipt: ',
			await transaction.receipt
		);

		return contract;
	})
	.then(async contract => {
		const response = await contract.checkGoalReached();

		console.log(response);
	})
	.catch(error => console.log(error));
