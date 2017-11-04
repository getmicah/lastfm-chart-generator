"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const CollageGenerator_1 = require("./CollageGenerator");
const process = require("process");
class Program {
    main(args) {
        const params = this.parseArgs(args);
        if (params) {
            const chart = new CollageGenerator_1.default(params[0], params[1], params[2]);
            chart.load();
        }
    }
    parseArgs(args) {
        if (args.length !== 5) {
            this.error('Invalid arguments');
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
                this.error('Invalid period');
                return;
        }
        const size = Number(args[4]);
        const validSize = size === 3 || size === 4 || size === 5;
        if (!validSize) {
            this.error('Invalid size');
            return;
        }
        return [user, period, size];
    }
    error(message) {
        console.error(`Error: ${message}`);
        console.log('Usage: lcg <user> <period> <size>');
        console.log('Params:');
        console.log('user\t<last.fm username>');
        console.log('period\tweek, month, 3month, 6month, year, overall');
        console.log('size\t3, 4, 5');
    }
}
new Program().main(process.argv);
//# sourceMappingURL=index.js.map