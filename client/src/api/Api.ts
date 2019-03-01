export class Api {

	// static url = 'http://127.0.0.1:21002/';
	static url = 'http://dinosaur-wechat-test.campfiregames.cn/wstail/';

	static async get(uri: string) {

		const x = await fetch(this.url + uri);

		const ab = await x.arrayBuffer();

		return new Uint8Array(ab);
	}
}
