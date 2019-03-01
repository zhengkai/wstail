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

	baseUrl = 'ws://dinosaur-wechat-test.campfiregames.cn/wstail/';

	connectID = 0;

	connectCount = 0;

	conn: IWsConn = null;

	fileName = '';

	ts: number;

	cb: any;

	stop: true;

	constructor() {
	}

	close() {
		this.stop = true;
		this._disconnect();
		this.conn = null;
	}

	_disconnect() {
		this.connectID++;
		if (this.conn !== null) {
			this.conn.ws.close();
		}
	}

	connect() {

		this.ts = Date.now();
		this._disconnect();

		let ws = new WebSocket(this.baseUrl + 'listen');
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

		if (this.stop) {
			return;
		}

		this.connectCount++;
		if (this.connectCount > 5) {
			this.connectCount = 5;
		}

		setTimeout(() => {
			this.connect();
		}, this.connectCount * 1000);
	}

	async onerror(e, id) {
		console.warn('ws onerror', e);
	}

	async onmessage(e, id) {

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

		this.connectCount = 0;

		const login = {
			fileName: this.fileName,
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
