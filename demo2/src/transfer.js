const evmlib = require('evm-lite-lib');

// Transaction Addresses
const from = '0x479e8b1b9d8b509755677f6d61d2f7339ba4c0fd';
const to = '0x1dEC6F07B50CFa047873A508a095be2552680874';

// EVMLC object
const evmlc = new evmlib.EVMLC('127.0.0.1', 8080, {
	from,
	gas: 100000,
	gasPrice: 0
});

// Data directory object
const dataDirectory = new evmlib.DataDirectory('[..]/.evmlc');

// Async Await
async function signTransactionLocallyAsync() {
	// Get keystore object from the keystore directory and decrypt
	const account = await dataDirectory.keystore.decryptAccount(
		from,
		'password'
	);

	// Prepare a transaction with value of 2000
	const transaction = await evmlc.prepareTransfer(to, 2000);

	// Sign transaction and return the same Transaction object
	// Send transaction to node
	await transaction.submit({}, account);

	return transaction;
}

signTransactionLocallyAsync()
	.then(transaction => console.log(transaction.hash))
	.catch(error => console.log(error));
