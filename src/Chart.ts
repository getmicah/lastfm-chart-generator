import { createCanvas, Image } from 'canvas';
import fetch from 'node-fetch';
import * as fs from 'fs';

interface album {
	name: string;
	playcount: number;
	mbid: string;
	url: string;
	artist: {
		name: string;
		mbid: string;
		url: string;
	}
	image: {
		"#text": string;
		size: string;
	}[];
	"@attr": {
		rank: number
	}
}

interface Canvas extends HTMLCanvasElement {
	toBuffer(): BufferSource
}

export default class Chart {
	user: string;
	period: string;
	size: number;
	constructor(user, period, size) {
		this.user = user;
		this.period = period;
		this.size = size;
	}
	public load(): void {
		const api = 'http://ws.audioscrobbler.com/2.0/';
		const key = '91289adaabc4ed1a559d5928015cd702';
		const method = 'user.gettopalbums';
		const limit = this.size * this.size;
		const req = `${api}?method=${method}&user=${this.user}&period=${this.period}&limit=${limit}&api_key=${key}&format=json`;
		fetch(req).then(r => r.json())
			.then(this.parse.bind(this))
			.then(this.draw.bind(this))
			.then(this.download.bind(this))
			.catch(this.error.bind(this));
	}
	private parse(r): Promise<Image[]> {
		return new Promise<Image[]>((resolve, reject) => {
			if (r.error) {
				reject(r.message);
			}
			const albums: album[] = r.topalbums.album;
			if (albums.length == 0) {
				reject('User has no scrobbles during the given period');
			}
			const covers: Image[] = [];
			let queue = albums.length;
			for (let i = 0; i < albums.length; i++) {
				let img: Image = new Image;
				fetch(albums[i].image[3]['#text'])
					.then((r) => r.buffer())
					.then((b) => {
						img.src = b;
						queue--;
						if (queue <= 0) {
							resolve(covers);
						}
					}).catch(() => reject());
				img.crossOrigin = 'anonymous';
				img.name = albums[i].name;
				img.artist = albums[i].artist.name;
				covers.push(img);
			}
		});
	}
	private draw(covers: Image[]): Canvas {
		const dx = covers[0].width;
		const w = dx * this.size;
		const canvas: Canvas = createCanvas(w, w);
		const ctx = canvas.getContext('2d');
		const pad = 3;
		const fontSize = 13;
		let i = 0;
		for (let y = 0; y < canvas.width; y+=dx) {
			for (let x = 0; x < canvas.height; x+=dx) {
				ctx.drawImage(covers[i], x, y);
				ctx.font = `${fontSize}px monospace`;
				ctx.fillStyle = 'black';
				ctx.fillText(`${covers[i].artist}`, x+pad, y+pad+10);
				ctx.fillText(`${covers[i].name}`, x+pad, y+pad+10+fontSize);
				ctx.fillStyle = 'white';
				ctx.fillText(`${covers[i].artist}`, x+pad+1, y+pad+9);
				ctx.fillText(`${covers[i].name}`, x+pad+1, y+pad+9+fontSize);
				i++;
			}
		}
		return canvas;
	}
	private download(canvas: Canvas) {
		const buf = canvas.toBuffer();
		fs.writeFileSync('chart.png', buf);
		console.log('created chart.png');
	}
	private error(msg: string): void {
		console.error(msg);
	}
}