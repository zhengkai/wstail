import { Injectable } from '@angular/core';

interface IWsConn {
	ws: WebSocket;
	id: number;
}

@Injectable({
	providedIn: 'root'
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
		console.log('onmessage', e);
	}

	onopen(e) {
		console.log('onopen', e);
	}
}
