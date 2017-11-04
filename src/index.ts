import Chart from './Chart';
import * as process from 'process';

class Program {
	public main(args: string[]): void {
		const params = this.parseArgs(args);
		if (params) {
			const chart = new Chart(params[0], params[1], params[2]);
			chart.load();
		}
	}
	private parseArgs(args: string[]): [string, string, number] {
		if (args.length !== 5) {
			this.usage('Invalid arguments');
			return;
		}
		const user = args[2];
		let period;
		switch (args[3]) {
			case 'week':
				period = '7day';
				break;
			case 'month':
				period = '1month';
				break;
			case '3month':
				period = '3month';
				break;
			case '6month':
				period = '6month';
				break;
			case 'year':
				period = '12month';
				break;
			case 'overall':
				period = 'overall';
				break;
			default:
				this.usage('Invalid period');
				return;
		}
		const size = Number(args[4]);
		const validSize = size === 3 || size === 4 || size === 5;
		if (!validSize) {
			this.usage('Invalid size');
			return;
		}
		return [user, period, size];
	}
	private usage(message: string): void {
		console.error(`Error: ${message}`);
		console.log('Usage: lcg <user> <period> <size>');
		console.log('Params:');
		console.log('user\t<last.fm username>');
		console.log('period\tweek, month, 3month, 6month, year, overall');
		console.log('size\t3, 4, 5');
		console.log('https://github.com/getmicah/lastfm-chart-generator');
	}
}

new Program().main(process.argv);