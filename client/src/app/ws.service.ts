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

	count = 0;

	conn: IWsConn = null;

	ts: number;

	cb: any;

	constructor() {
		this.connect();
	}

	connect() {

		this.connectID++;
		this.ts = Date.now();

		if (this.conn !== null) {
			console.log('conn not null', this.conn);
			this.conn.ws.close();
		}

		let ws = new WebSocket('ws://127.0.0.1:21002/ws/listen');
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
			this.onopen(e, id);
		};
		ws.onclose = (e) => {
			if (id !== this.connectID) {
				return;
			}
			this.connectID++;
			this.onclose(e, id);
			ws = null;
		};
		ws.onmessage = (e) => {
			if (id !== this.connectID) {
				return;
			}
			this.onmessage(e, id);
		};
		ws.onerror = (e) => {
			if (id !== this.connectID) {
				return;
			}
			this.onerror(e, id);
		};
	}

	async onclose(e, id) {
		console.log('onclose', e.code, e.reason, Date.now() - this.ts, e);
	}

	async onerror(e, id) {
		console.log('onerror', e);
	}

	async onmessage(e, id) {
		// console.log('onmessage', e);
		//

		const ab = await (new Response(e.data)).arrayBuffer();

		const x = await pb.MsgReturn.decode(new Uint8Array(ab));

		for (const a of x.msg) {

			let v: any;
			const t = a.type_url.substring(lengthPrefixType);

			switch (t) {
			case 'Update':
				v = await pb.Update.decode(a.value);
				break;
			case 'PrevContent':
				v = await pb.PrevContent.decode(a.value);
				break;
			}

			console.log(x, ab.byteLength, a.value.length, lengthPrefixType);

			if (this.cb) {
				this.cb.recv.call(this.cb, this, v, t, id);
			}
		}
	}

	async onopen(e, id) {
		console.log('onopen', e);

		const login = {
			name: 'rpg',
			connectType: pb.Login.ConnectType.NEW,
		} as pb.Login;

		const ab = pb.Login.encode(login).finish();
		this.send(id, ab);

		if (this.cb) {
			this.cb.connect.call(this.cb, this, id);
		}
	}

	async send(id: number, message: Uint8Array|ArrayBuffer) {
		if (id !== this.conn.id) {
			return;
		}
		this.conn.ws.send(message);
	}
}
