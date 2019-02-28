import { AfterViewChecked, Component, OnInit, ElementRef, Renderer2, ViewChild } from '@angular/core';
import { WSService } from '../ws.service';
import { Api } from '../../api';
import { pb } from '../../pb/pb';

@Component({
	selector: 'app-view',
	templateUrl: './view.component.html',
	styleUrls: ['./view.component.scss'],
})
export class ViewComponent implements OnInit, AfterViewChecked {

	@ViewChild('scroll') private sc: ElementRef;

	msgPool: Array<string> = [];

	file: Array<string> = [];
	fileSelect = '';

	charNum = 0;

	select(fileName: string) {
		this.ws.cb = this;
		this.fileSelect = fileName;
		this.ws.fileName = fileName;
		this.ws.connect();
		this.msgPool.length = 0;
	}

	async fetch() {

		const form = new FormData();
		form.append('type', 'test');
		form.append('v', 'yes rpg');

		const x = await fetch('https://dinosaur-wechat-test.campfiregames.cn/logjson', {
			method: 'POST',
			body: form,
		});
		console.log('fetch', x.status, x.statusText);
	}

	async list() {

		const ab = await Api.get('file');

		const r = pb.FileReturn.decode(ab);

		this.file.length = 0;

		r.file.forEach((s) => {
			this.file.push(s);
		});

		console.log(this.file);
	}

	constructor(private ws: WSService, private renderer: Renderer2) {
		this.fetch();
		this.list();
	}

	recv(ws: WSService, msg: any, t: string, id: number) {

		if (t === 'Update' && msg.reset) {
			this.msgPool.length = 0;
			this.charNum = 0;
		}

		while (true) {
			if (this.charNum < 10000000) {
				break;
			}
			const s = this.msgPool.shift();
			this.charNum -= s.length;
		}

		const content = msg.msg;
		this.charNum += content.length;

		switch (t) {
		case 'Update':
			this.msgPool.push(content);
			break;
		case 'PrevContent':
			this.msgPool.unshift(content);
			break;
		}
	}

	ngAfterViewChecked() {
		this.scrollDown();
	}

	scrollDown() {
		const el = this.sc.nativeElement;
		this.renderer.setProperty(
			el,
			'scrollTop',
			Number.MAX_SAFE_INTEGER,
		);
	}

	connect(ws: WSService, id: number) {
		console.log('connect', id);
	}

	ngOnInit() {
	}
}
