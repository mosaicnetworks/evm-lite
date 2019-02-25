const { EVMLC } = require('evm-lite-lib');
const { DataDirectory } = require('evm-lite-lib');

const solc = require('solc');
const fs = require('fs');

// Default from address
const from = '0x1a5c6b111e883d920fd24fee0bafae838958fa05';

// EVMLC object
const evmlc = new EVMLC('127.0.0.1', 8080, {
	from,
	gas: 1000000,
	gasPrice: 0
});

// Keystore object
const dataDirectory = new DataDirectory('[..]/.evmlc');

// Contract Object
const contractPath = '../assets/CrowdFunding.sol';
const contractFile = fs.readFileSync(contractPath, 'utf8');
const contractName = ':' + 'CrowdFunding';

const output = solc.compile(contractFile, 1);
const ABI = JSON.parse(output.contracts[contractName].interface);
const data = output.contracts[contractName].bytecode;
const account = dataDirectory.keystore.decryptAccount(from, 'password');

const loadContract = async () => {
	// Generate contract object with ABI and data
	const contract = await evmlc.loadContract(ABI, {
		data
		// Will generate functions for the deployed contract at the address if set.
		/* contractAddress: '' */
	});

	// Deploy and return contract with functions populated
	return await contract.deploy(await account, [10000], {
		// by default
		gas: evmlc.defaultGas,
		gasPrice: evmlc.defaultGasPrice
	});
};

loadContract()
	.then(async contract => {
		const transaction = await contract.methods.contribute();

		transaction.value(200);

		await transaction.sign(await account);
		await transaction.submit();

		return contract;
	})
	.then(async contract => {
		const account = await evmlc.getAccount(contract.options.address.value);
		console.log(account);

		return contract;
	})
	.then(async contract => {
		const transaction = await contract.methods.checkGoalReached();

		await transaction.sign(await account);

		const response = await transaction.submit();
		console.log(response);

		return contract;
	})
	.catch(error => console.log(error));
