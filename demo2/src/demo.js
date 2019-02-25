util = require("util");
JSONbig = require('json-bigint');
argv = require('minimist')(process.argv.slice(2));

evml = require('evm-lite-lib');
const Keystore = evml.Keystore;

function DemoNode(name, host, port) {
    this.name = name
    this.api = new evml.EVMLC(host, port, {})
    this.account = {}
}

var _keystore;
var _accounts = [];
var _nodes = [];


init = async function() {
    console.log(argv);

    var ips = argv.ips.split(",").sort();
    console.log("sorted ips: ", ips);

    var port = argv.port;
    
    _keystore = new Keystore(argv.keystore, "keystore")

    for (i=0; i<ips.length; i++) {
        node = new DemoNode(
            util.format('node%d', i+1),
            ips[i],
            port);
        _nodes.push(node);
    }

    return
}

decryptAccounts = async function() {
    accounts = await _keystore.list()

    for (const element of accounts ) {
        acc = await _keystore.decryptAccount(element.address, "supersecurepassword")
        _accounts.push(acc);
    }

    console.log("accounts", _accounts);

    //assuming that there are as many nodes as accounts
    for (i=0; i<_nodes.length; i++) {
        _nodes[i].account = _accounts[i];
    }

    console.log("nodes", _nodes);
}

displayAllBalances = async function() {
    for (const node of _nodes) {
        acc = await node.api.accounts.getAccount(node.account.address)
        console.log(node.name, acc)
    }
}

transferRaw = async function(fromNode, to, amount) {    
    const tx = await fromNode.api.accounts.prepareTransfer(to, amount);
    await tx.submit({}, fromNode.account);
    return tx;
    
}

/******************************************************************************/

init()
.then(() => decryptAccounts())
.then(() => displayAllBalances())
.then(() => transferRaw(_nodes[0], _nodes[1].account.account.address, 666))
.then(() => displayAllBalances())