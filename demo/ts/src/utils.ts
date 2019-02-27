import * as prompt from 'prompt';
import * as util from 'util';

export const FgRed = '\x1b[31m';
export const FgGreen = '\x1b[32m';
export const FgYellow = '\x1b[33m';
export const FgBlue = '\x1b[34m';
export const FgMagenta = '\x1b[35m';
export const FgCyan = '\x1b[36m';
export const FgWhite = '\x1b[37m';

export const log = (color: string, text: string) => {
	console.log(color + text + '\x1b[0m');
};

export const step = (message: string) => {
	log(FgWhite, '\n' + message);
	return new Promise(resolve => {
		prompt.get('PRESS ENTER TO CONTINUE', (err: any, res: any) => {
			resolve();
		});
	});
};

export const explain = (message: string) => {
	log(FgCyan, util.format('\nEXPLANATION:\n%s', message));
};

export const space = () => {
	console.log('\n');
};

export const sleep = (time: number) => {
	return new Promise(resolve => setTimeout(resolve, time));
};
