// const decryptAccounts = async function() {
// 	accounts = await _keystore.list();

// 	for (const element of accounts) {
// 		acc = await _keystore.decryptAccount(
// 			element.address,
// 			'supersecurepassword'
// 		);
// 		_accounts.push(acc);
// 	}

// 	console.log('accounts', _accounts);

// 	//assuming that there are as many nodes as accounts
// 	for (i = 0; i < _nodes.length; i++) {
// 		_nodes[i].api.defaultFrom = _accounts[i].address;
// 		_nodes[i].account = _accounts[i];
// 	}

// 	console.log('nodes', _nodes);
// };

// const displayAllBalances = async function() {
// 	for (const node of _nodes) {
// 		acc = await node.api.accounts.getAccount(node.account.address);
// 		console.log(node.name, acc);
// 	}
// };

// const transferRaw = async function(fromNode, to, amount) {
// 	const tx = await fromNode.api.accounts.prepareTransfer(to, amount);

// 	await tx.submit({}, fromNode.account);

// 	return tx;
// };

// const compileAndCreateContract = () => {
// 	// pass
// };

// class CrowdFunding {
// 	contract;
// 	account;

// 	constructor(contract, account) {
// 		this.contract = contract;
// 		this.account = account;
// 	}

// 	async deploy() {
// 		// Deploy crowdfunding contract with constructor param as 1000
// 		const transaction = await this.contract.deploy(account, [1000]);

// 		return transaction;
// 	}

// 	async contribute() {
// 		const transaction = await this.contract.methods.contribute(1000);

// 		await transaction.submit({}, this.account);

// 		return transaction;
// 	}

// 	async checkGoalReached() {
// 		const transaction = await this.contract.methods.checkGoalReached();
// 		const response = await transaction.submit({}, this.account);

// 		return response;
// 	}
// }

// /******************************************************************************/

// init()
// 	.then(() => decryptAccounts())
// 	.then(() => displayAllBalances())
// 	.then(() => transferRaw(_nodes[0], _nodes[1].account.account.address, 666))
// 	.then(() => displayAllBalances());

// const oneTimeTransfer = async () => {
//     const keystore = new evmlc.Keystore(
//         '/Users/danu/Desktop/evm-lite/demo2/src/',
//         'keystore'
//     );
//     const evmlcNode = new evmlc.EVMLC('127.0.0.1', 8080, {
//         from: '0xA4a5F65Fb3752b2B6632F2729f17dd61B2aaD650',
//         gas: 100000,
//         gasPrice: 0
//     });
//     const account = await keystore.decryptAccount(
//         evmlcNode.defaultFrom,
//         'supersecurepassword'
//     );
//     const others = [
//         '0xfFFFC2A95F453aB0BC474c04eFceCf7cB47A29Aa',
//         '0xbB11ae377c9a20bf12F322048B71864f2911e476',
//         '0x2Fa4d156f0Ac83C792B7d948983A6c60957c1779'
//     ];

//     for (const other of others) {
//         const transaction = await evmlcNode.accounts.prepareTransfer(
//             other,
//             100000000
//         );
//         await transaction.submit({}, account);

//         console.log(await transaction.receipt);
//     }
// };
