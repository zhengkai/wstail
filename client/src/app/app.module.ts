import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { RoutingModule } from './routing/routing.module';
import { BootstrapComponent } from './bootstrap/bootstrap.component';
import { ViewComponent } from './view/view.component';

@NgModule({
	declarations: [
		BootstrapComponent,
		ViewComponent,
	],
	imports: [
		BrowserModule,
		RoutingModule,
	],
	providers: [],
	bootstrap: [
		BootstrapComponent,
	],
})
export class AppModule {}
