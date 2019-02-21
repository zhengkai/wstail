import { AfterViewChecked, Component, OnInit, ElementRef, Renderer2, ViewChild } from '@angular/core';
import { WSService } from '../ws.service';

@Component({
	selector: 'app-view',
	templateUrl: './view.component.html',
	styleUrls: ['./view.component.scss'],
})
export class ViewComponent implements OnInit, AfterViewChecked {

	@ViewChild('scroll') private sc: ElementRef;

	msgPool: Array<string> = [];

	charNum = 0;

	constructor(private ws: WSService, private renderer: Renderer2) {
		ws.cb = this;
	}

	recv(ws: WSService, msg: string, id: number) {
		this.charNum += msg.length;
		while (true) {
			if (this.charNum < 10000000) {
				break;
			}
			const s = this.msgPool.shift();
			this.charNum -= s.length;
		}
		this.msgPool.push(msg);
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
