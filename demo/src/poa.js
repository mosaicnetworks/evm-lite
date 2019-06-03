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

const online = true;

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
	this.ip = host;
}

const allAccounts = [];
const allNodes = [];

var crowdFunding = {};
var contractPath = '';
var nodeno = 0;


const init = async () => {
	console.group('Initialize Nodes: ');

	const ips = argv.ips 
		.replace(/\s/g, '')
		.split(',');
//		.sort();
	const port = argv.port;
	const keystore = new evmlc.Keystore(argv.keystore, 'keystore');
	const passwordPath = argv.pwd;
	nodeno = argv.nodeno;
	const password = readPasswordFile(passwordPath);

	contractPath = argv.contract;

	console.log('Sorted IPS: ', ips);
    console.log('Node No: ', nodeno);
	console.log('Keystore Path: ', keystore.path);
	console.log('Password File Path: ', passwordPath);
	console.log('Contract File Path: ', contractPath);

	for (i = 0; i < ips.length; i++) {
		node = new Node(util.format('node%s',nodeno ), ips[i], port);
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

const compiledSmartContract = async () => {
	const input = fs.readFileSync(contractPath, {
		encoding: 'utf8'
	});
	const output = solc.compile(input.toString(), 1);
	const bytecode = output.contracts[`:POA_Genesis`].bytecode;
	const abi = output.contracts[`:POA_Genesis`].interface;

	const contract = await allNodes[0].api.contracts.load(JSON.parse(abi), {
		data: bytecode
	});

	return contract;
};

class GenesisContract {
	constructor (account, ip) {
		this.account = account;
		this.contract = new evmlc.Contract({
            "gas":100000000, 
            "gasPrice": 0, 
            "from": this.account.address, 
            "address": "0xabbaabbaabbaabbaabbaabbaabbaabbaabbaabba", 
            "interface": [{"constant":true,"inputs":[{"name":"_publicKey","type":"bytes32"}],"name":"checkAuthorisedPublicKey","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_address","type":"address"}],"name":"dev_isGenesisWhitelisted","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[],"name":"dev_27","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"_address","type":"address"}],"name":"checkAuthorised","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"whiteList","outputs":[{"name":"person","type":"address"},{"name":"flags","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"dev_getSender","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"dev_getWhitelistCount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_nomineeAddress","type":"address"},{"name":"_accepted","type":"bool"}],"name":"castNomineeVote","outputs":[{"name":"decided","type":"bool"},{"name":"voteresult","type":"bool"}],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[],"name":"dev_getGenesisWhitelist0","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"_address","type":"address"}],"name":"dev_isWhitelisted","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_address","type":"address"}],"name":"dev_getCurrentNomineeVotes","outputs":[{"name":"yes","type":"uint256"},{"name":"no","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_nomineeAddress","type":"address"},{"name":"_moniker","type":"bytes32"}],"name":"submitNominee","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"inputs":[{"name":"_moniker","type":"bytes32"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_nominee","type":"address"},{"indexed":false,"name":"_yesVotes","type":"uint256"},{"indexed":false,"name":"_noVotes","type":"uint256"},{"indexed":true,"name":"_accepted","type":"bool"}],"name":"NomineeDecision","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_nominee","type":"address"},{"indexed":true,"name":"_voter","type":"address"},{"indexed":false,"name":"_yesVotes","type":"uint256"},{"indexed":false,"name":"_noVotes","type":"uint256"},{"indexed":true,"name":"_accepted","type":"bool"}],"name":"NomineeVoteCast","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_nominee","type":"address"},{"indexed":true,"name":"_proposer","type":"address"}],"name":"NomineeProposed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_address","type":"address"},{"indexed":true,"name":"_moniker","type":"bytes32"}],"name":"MonikerAnnounce","type":"event"}]
            },  ip , 8080 );
    }

    async genericTrans(functionname, value, params) {
        console.log(functionname, params);

        const transaction = await this.contract.methods[functionname](...params);
    
        transaction.value(value);

        console.log("From", transaction.tx.from) ;
        console.log("Data", transaction.tx.data) ;

        const response = await transaction.submit(this.account, { timeout: 3 });
  
        if (value > 0)
        {
           const receipt = await transaction.receipt;
  
           console.dir(receipt, { depth: 6, colors: true });
  
           const logs = this.contract.parseLogs(receipt.logs);
  
           for (const log of logs) {
               console.log(
                    log.event || 'No Event Name',
                    JSON.stringify(log.args, null, 2)
               );
           }
        }
        else
        {
           console.dir(response, { depth: 6, colors: true });
        }
        return transaction;
     }
}
	