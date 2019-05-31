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
		this.contract = new evmlc.Contract({"gas":100000000, "gasPrice": 0, 
"from": this.account.address, 
"address": "0xabbaabbaabbaabbaabbaabbaabbaabbaabbaabba", "interface": 
[{"constant":true,"inputs":[{"name":"_publicKey","type":"bytes32"}],"name":"checkAuthorisedPublicKey","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_address","type":"address"}],"name":"dev_isGenesisWhitelisted","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[],"name":"dev_27","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"_address","type":"address"}],"name":"checkAuthorised","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"whiteList","outputs":[{"name":"person","type":"address"},{"name":"flags","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"dev_getSender","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"dev_getWhitelistCount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_nomineeAddress","type":"address"},{"name":"_accepted","type":"bool"}],"name":"castNomineeVote","outputs":[{"name":"decided","type":"bool"},{"name":"voteresult","type":"bool"}],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[],"name":"dev_getGenesisWhitelist0","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"pure","type":"function"},{"constant":true,"inputs":[{"name":"_address","type":"address"}],"name":"dev_isWhitelisted","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"_address","type":"address"}],"name":"dev_getCurrentNomineeVotes","outputs":[{"name":"yes","type":"uint256"},{"name":"no","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_nomineeAddress","type":"address"},{"name":"_moniker","type":"bytes32"}],"name":"submitNominee","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"inputs":[{"name":"_moniker","type":"bytes32"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_nominee","type":"address"},{"indexed":false,"name":"_yesVotes","type":"uint256"},{"indexed":false,"name":"_noVotes","type":"uint256"},{"indexed":true,"name":"_accepted","type":"bool"}],"name":"NomineeDecision","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_nominee","type":"address"},{"indexed":true,"name":"_voter","type":"address"},{"indexed":false,"name":"_yesVotes","type":"uint256"},{"indexed":false,"name":"_noVotes","type":"uint256"},{"indexed":true,"name":"_accepted","type":"bool"}],"name":"NomineeVoteCast","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_nominee","type":"address"},{"indexed":true,"name":"_proposer","type":"address"}],"name":"NomineeProposed","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"_address","type":"address"},{"indexed":true,"name":"_moniker","type":"bytes32"}],"name":"MonikerAnnounce","type":"event"}]
},  ip , 8080 );
	}
	


async genericTrans(functionname, value, params) {

      console.log(functionname, params);
  //    console.log('Contract: ', this.contract);
      const transaction = await this.contract.methods[functionname](...params);
//		console.log('Transaction Created');

      transaction.value(value);
//      console.log("Value Set");
//      console.dir(transaction, { depth: 6, colors: true });
      console.log("From", transaction.tx.from) ;
      console.log("Data", transaction.tx.data) ;
//      console.log('Transaction: ', transaction.parse(), '\n');
      const response = await transaction.submit(this.account, { timeout: 3 });

//      console.dir(response, { depth: 6, colors: true });
    //  console.log('Response: ', response, '\n');
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


init()
        .then(object => decryptAccounts(object))
        .then(() => step('STEP 0.1) Get ETH Account on this Node'))
        .then(() => {
                space();
                return displayAllBalances();
        })
        .then(() =>
                explain(
                        'Quick check that we can talk to EVM-Lite \n'
                )
        )
        .then(() => step('STEP 0.2) Set Up Object for accessing for Genesis Authority Smart Contract'))

        .then(async contract => {
                genesis = new GenesisContract(allNodes[nodeno].account, allNodes[nodeno].ip);
                await genesis.genericTrans('dev_27', 0, []);
                await genesis.genericTrans('dev_getGenesisWhitelist0', 0, []);

	})
        .then(() =>
                explain(
                        'Quick check that things are working. We should have a return value of 27 from the test function and the first hard coded Genesis whitelist entry.  \n'
                )
        )


        .then(() => step('STEP 1) Am I on the whitelist or genesis whitelist?'))

        .then(async contract => {
                await genesis.genericTrans('dev_isWhitelisted', 0, [allNodes[nodeno].account.address]);
                await genesis.genericTrans('dev_isGenesisWhitelisted', 0, [allNodes[nodeno].account.address]);
                await genesis.genericTrans('dev_getWhitelistCount', 0, []);

                await genesis.genericTrans('checkAuthorised', 0, [allNodes[nodeno].account.address]);
	})
        .then(() =>
                explain(
                        'Check the authority of this account.  \n'
                )
        )




	.then(() => step('STEP 2) Node 0 Nominates Node 2'))
	.then(async contract => {    // submitNominee (address _nomineeAddress, bytes32 _moniker) public payable checkAuthorisedModifier(msg.sender)

		if ( nodeno == 0)
		{ 	
                    await genesis.genericTrans('submitNominee', 10000000, [allNodes[2].account.address, 'Node 2']);
                    await genesis.genericTrans('dev_getCurrentNomineeVotes', 0, [allNodes[2].account.address]);
                    await genesis.genericTrans('dev_getWhitelistCount', 0, []);
                    await genesis.genericTrans('dev_isWhitelisted', 0, [allNodes[nodeno].account.address]);

		    
		}
		else
		{
		    console.log(FgGreen, 'Do nothing this step');	
		}

	})

        .then(() =>
                explain(
                        ' \n'
                )
        )
	.then(() => {process.exit(0);} )


        .then(() => step('STEP 3) Node 0 Votes for Node 2'))
	.then(async contract => {    // function castNomineeVote(address _nomineeAddress, bool _accepted) public payable checkAuthorisedModifier(msg.sender) returns (bool decided, bool voteresult)

		if ( nodeno == 0)
		{ 	
                    await genesis.genericTrans('castNomineeVote', 10000000, [allNodes[2].account.address, true]);
                    await genesis.genericTrans('dev_getCurrentNomineeVotes', 0, [allNodes[2].account.address]);
		}
		else
		{
		    console.log(FgGreen, 'Do nothing this step');	
		}

	})

        .then(() =>
                explain(
                        ' \n'
                )
        )
        .then(() => step('STEP 4) Node 2 Tries to Join'))
	.then(async contract => {    // submitNominee (address _nomineeAddress, bytes32 _moniker) public payable checkAuthorisedModifier(msg.sender)

		if ( nodeno == 2)
		{ 	
		    console.log(FgYellow, 'Restart Node 2');
		}
		else
		{
		    console.log(FgGreen, 'Do nothing this step');	
		}

	})

        .then(() =>
                explain(
                        ' \n'
                )
        )
        .then(() => step('STEP 5) Node 1 Votes for Node 2'))
	.then(async contract => {    // function castNomineeVote(address _nomineeAddress, bool _accepted) public payable checkAuthorisedModifier(msg.sender) returns (bool decided, bool voteresult)

		if ( nodeno == 1)
		{ 	
                    await genesis.genericTrans('castNomineeVote', 100000, [allNodes[2].account.address,true]);
                    await genesis.genericTrans('dev_getCurrentNomineeVotes', 0, [allNodes[2].account.address]);
		}
		else
		{
		    console.log(FgGreen, 'Do nothing this step');	
		}

	})

        .then(() =>
                explain(
                        ' \n'
                )
        )
        .then(() => step('STEP 6) Node 2 joins'))
        .then(() =>
                explain(
                        ' \n'
                )
        )
        .then(() => step('STEP 7) Node 0 Nominates Node 3'))
	.then(async contract => {    // submitNominee (address _nomineeAddress, bytes32 _moniker) public payable checkAuthorisedModifier(msg.sender)

		if ( nodeno == 0)
		{ 	
                    await genesis.genericTrans('submitNominee', 100000, [allNodes[3].account.address, 'Node 3']);
		}
		else
		{
		    console.log(FgGreen, 'Do nothing this step');	
		}

	})

        .then(() =>
                explain(
                        ' \n'
                )
        )
        .then(() => step('STEP 8) Node 0 Votes for Node 3'))
	.then(async contract => {    // function castNomineeVote(address _nomineeAddress, bool _accepted) public payable checkAuthorisedModifier(msg.sender) returns (bool decided, bool voteresult)

		if ( nodeno == 0)
		{ 	
                    await genesis.genericTrans('castNomineeVote', 100000, [allNodes[3].account.address, true]);
		}
		else
		{
		    console.log(FgGreen, 'Do nothing this step');	
		}

	})

        .then(() =>
                explain(
                        ' \n'
                )
        )
        .then(() => step('STEP 9) Node 1 Votes for Node 3'))
	.then(async contract => {    // function castNomineeVote(address _nomineeAddress, bool _accepted) public payable checkAuthorisedModifier(msg.sender) returns (bool decided, bool voteresult)

		if ( nodeno == 1)
		{ 	
                    await genesis.genericTrans('castNomineeVote', 100000, [allNodes[3].account.address, true]);
		}
		else
		{
		    console.log(FgGreen, 'Do nothing this step');	
		}

	})

        .then(() =>
                explain(
                        ' \n'
                )
        )
        .then(() => step('STEP 10) Node 2 Votes against Node 3'))
	.then(async contract => {    //function castNomineeVote(address _nomineeAddress, bool _accepted) public payable checkAuthorisedModifier(msg.sender) returns (bool decided, bool voteresult)

		if ( nodeno == 2)
		{ 	
                    await genesis.genericTrans('castNomineeVote', 100000, [allNodes[3].account.address, false]);
		}
		else
		{
		    console.log(FgGreen, 'Do nothing this step');	
		}

	})

        .then(() =>
                explain(
                        ' \n'
                )
        )
        .then(() => step('STEP 11) Node 3 tries to join'))
	.then(async contract => {    // submitNominee (address _nomineeAddress, bytes32 _moniker) public payable checkAuthorisedModifier(msg.sender)

		if ( nodeno == 0)
		{ 	
                    console.log(FgYellow, "Node 3 joins");
		}
		else
		{
		    console.log(FgGreen, 'Do nothing this step');	
		}

	})

        .then(() =>
                explain(
                        ' \n'
                )
        )



        .catch(err => log(FgRed, err));




