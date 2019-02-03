import { Injectable } from '@angular/core';
import { pb } from '../pb/pb';

interface IWsConn {
	ws: WebSocket;
	id: number;
}

const lengthPrefixType = 'type.googleapis.com/pb.'.length;

@Injectable({
	providedIn: 'root',
})
export class WSService {

	connectID = 0;

	conn: IWsConn = null;

	constructor() {
		this.connect();
	}

	connect() {

		this.connectID++;

		if (this.conn !== null) {
			console.log('conn not null', this.conn);
			this.conn.ws.close();
		}

		const ws = new WebSocket('ws://127.0.0.1:21002/ws/listen');
		const id = this.connectID;
		const conn = {
			id,
			ws,
		} as IWsConn;

		this.conn = conn;

		ws.onopen = (e) => {
			if (id !== this.connectID) {
				return;
			}
			this.onopen(e);
		};
		ws.onclose = (e) => {
			if (id !== this.connectID) {
				return;
			}
			this.onclose(e);
		};
		ws.onmessage = (e) => {
			if (id !== this.connectID) {
				return;
			}
			this.onmessage(e);
		};
		ws.onerror = (e) => {
			if (id !== this.connectID) {
				return;
			}
			this.onerror(e);
		};
	}

	onclose(e) {
		console.log('onclose', e);
	}

	onerror(e) {
		console.log('onerror', e);
	}

	onmessage(e) {
		console.log('onmessage', e, lengthPrefixType);

		(async () => {
			// const buffer = await e.data();

			const ab = await (new Response(e.data)).arrayBuffer();

			const r = pb.MsgA.decode(new Uint8Array(ab));

			for (const o of r.msg) {

				// console.log(o);

				let fn: any;

				switch (o.type_url.substring(lengthPrefixType)) {
				case 'Play':
					fn = pb.Play;
					break;
				case 'GameAuth':
					fn = pb.GameAuth;
					break;
				}

				const a = fn.decode(o.value);
				console.log(a);
			}

			// console.log(r);
		})();
	}

	onopen(e) {
		console.log('onopen', e);
	}
}
