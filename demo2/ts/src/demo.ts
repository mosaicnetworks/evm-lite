import * as fs from 'fs';
import * as solc from 'solc';

import * as Utils from './utils';

import {
	Account,
	BaseContractSchema,
	EVMLC,
	Keystore,
	SolidityContract,
	Transaction
} from 'evm-lite-lib';

const defaultGas = 100000;
const defaultGasPrice = 0;

class Node {
	public readonly api: EVMLC;

	public account?: Account;

	constructor(public readonly name: string, host: string, port: number) {
		this.api = new EVMLC(host, port, {
			from: '',
			gas: defaultGas,
			gasPrice: defaultGasPrice
		});
	}
}

interface CrowdFundingSchema extends BaseContractSchema {
	contribute: () => Promise<Transaction>;
	checkGoalReached: () => Promise<Transaction>;
	settle: () => Promise<Transaction>;
}

class CrowdFunding {
	constructor(
		public readonly contract: SolidityContract<CrowdFundingSchema>,
		public readonly account: Account
	) {}

	public async deploy(value: number) {
		await this.contract.deploy(this.account, [value], { timeout: 2 });
		console.log('Receipt:', this.contract.receipt);
		return this;
	}

	public async contribute(value: number) {
		const transaction = await this.contract.methods.contribute();
		transaction.value(value);

		console.log('Transaction: ', transaction.parse(), '\n');
		await transaction.submit({ timeout: 2 }, this.account);

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

	public async checkGoalReached() {
		const transaction = await this.contract.methods.checkGoalReached();
		const response: any = await transaction.submit(
			{ timeout: 2 },
			this.account
		);

		const parsedResponse = {
			goalReached: response[0],
			beneficiary: response[1],
			fundingTarget: response[2].toFormat(0),
			current: response[3].toFormat(0)
		};

		Utils.log(Utils.FgBlue, JSON.stringify(parsedResponse, null, 2));

		return response;
	}

	public async settle() {
		const transaction = await this.contract.methods.settle();

		console.log('Transaction: ', transaction.parse(), '\n');
		await transaction.submit({ timeout: 2 }, this.account);

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

class Demo {
	public readonly accounts: Account[] = [];
	public readonly nodes: Node[] = [];

	public password: string;
	public contract: string;
	public keystore: Keystore;

	constructor(
		ips: string[],
		port: number,
		keystorePath: string,
		passwordPath: string,
		contractPath: string
	) {
		this.password = this.readFile(passwordPath);
		this.contract = this.readFile(contractPath);

		for (let i = 0; i < ips.length; i++) {
			const node = new Node(`node${i + 1}`, ips[i], port);

			this.nodes.push(node);
		}

		const data = this.parentAndName(keystorePath);
		this.keystore = new Keystore(data.parent, data.name);
	}

	public async decryptAccounts(): Promise<void> {
		const baseAccounts = await this.keystore.list(true, this.nodes[0].api);

		for (const baseAccount of baseAccounts) {
			const account = await this.keystore.decryptAccount(
				baseAccount.address,
				this.password
			);
			account.balance = baseAccount.balance;

			this.accounts.push(account);
		}

		this.setDefaultAccounts();
	}

	public async displayBalances(): Promise<void> {
		for (const node of this.nodes) {
			const baseAccount = await node.api.accounts.getAccount(
				node.account!.address
			);

			console.log(`${node.name}: `, '\n', baseAccount, '\n');
		}
	}

	public async transferRaw(
		from: Node,
		to: Node,
		value: number
	): Promise<void> {
		const transaction = await from.api.accounts.prepareTransfer(
			to.account!.address,
			value
		);

		await transaction.submit({ timeout: 2 }, from.account);
	}

	public async getContract(
		node: Node
	): Promise<SolidityContract<CrowdFundingSchema>> {
		const output = solc.compile(this.contract.toString(), 1);
		const bytecode = output.contracts[`:CrowdFunding`].bytecode;
		const abi = output.contracts[`:CrowdFunding`].interface;

		const contract = await node.api.contracts.load<CrowdFundingSchema>(
			JSON.parse(abi),
			{
				data: bytecode
			}
		);

		return contract;
	}

	private setDefaultAccounts(): void {
		this.nodes.forEach((node, index) => {
			node.api.defaultFrom = this.accounts[index].address;
			node.account = this.accounts[index];
		});
	}

	private readFile(path: string): string {
		return fs.readFileSync(path, { encoding: 'utf8' });
	}

	private parentAndName(path: string): { parent: string; name: string } {
		const list = path.split('/');
		const name = list.pop() || 'keystore';

		return {
			parent: list.join('/'),
			name
		};
	}
}
