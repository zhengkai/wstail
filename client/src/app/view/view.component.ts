import { Component, OnInit } from '@angular/core';
import { WSService } from '../ws.service';

@Component({
	selector: 'app-view',
	templateUrl: './view.component.html',
	styleUrls: ['./view.component.scss'],
})
export class ViewComponent implements OnInit {

	constructor(private ws: WSService) {
	}

	ngOnInit() {
	}
}
