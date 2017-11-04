"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const canvas_1 = require("canvas");
const node_fetch_1 = require("node-fetch");
const fs = require("fs");
class CollageGenerator {
    constructor(user, period, size) {
        this.user = user;
        this.period = period;
        this.size = size;
    }
    load() {
        const api = 'http://ws.audioscrobbler.com/2.0/';
        const key = '91289adaabc4ed1a559d5928015cd702';
        const method = 'user.gettopalbums';
        const limit = this.size * this.size;
        const req = `${api}?method=${method}&user=${this.user}&period=${this.period}&limit=${limit}&api_key=${key}&format=json`;
        node_fetch_1.default(req).then(r => r.json())
            .then(this.parse.bind(this))
            .then(this.draw.bind(this))
            .then(this.download.bind(this))
            .catch(this.error.bind(this));
    }
    parse(r) {
        return new Promise((resolve, reject) => {
            if (r.error) {
                reject(r.message);
            }
            const albums = r.topalbums.album;
            if (albums.length == 0) {
                reject('User has no scrobbles during the given period');
            }
            const covers = [];
            let queue = albums.length;
            for (let i = 0; i < albums.length; i++) {
                let img = new canvas_1.Image;
                node_fetch_1.default(albums[i].image[3]['#text'])
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
    draw(covers) {
        const dx = covers[0].width;
        const w = dx * this.size;
        const canvas = canvas_1.createCanvas(w, w);
        const ctx = canvas.getContext('2d');
        const pad = 3;
        const fontSize = 13;
        let i = 0;
        for (let y = 0; y < canvas.width; y += dx) {
            for (let x = 0; x < canvas.height; x += dx) {
                ctx.drawImage(covers[i], x, y);
                ctx.font = `${fontSize}px monospace`;
                ctx.fillStyle = 'black';
                ctx.fillText(`${covers[i].artist}`, x + pad, y + pad + 10);
                ctx.fillText(`${covers[i].name}`, x + pad, y + pad + 10 + fontSize);
                ctx.fillStyle = 'white';
                ctx.fillText(`${covers[i].artist}`, x + pad + 1, y + pad + 9);
                ctx.fillText(`${covers[i].name}`, x + pad + 1, y + pad + 9 + fontSize);
                i++;
            }
        }
        return canvas;
    }
    download(canvas) {
        const buf = canvas.toBuffer();
        fs.writeFileSync('collage.png', buf);
    }
    error(msg) {
        console.error(msg);
    }
}
exports.default = CollageGenerator;
//# sourceMappingURL=CollageGenerator.js.map