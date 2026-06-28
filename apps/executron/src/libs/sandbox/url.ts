export const urlPolyfill = `
class URLSearchParams {
	constructor(init) {
		this._params = [];
		if (typeof init === 'string') {
			if (init.startsWith('?')) init = init.slice(1);
			if (init) {
				for (const pair of init.split('&')) {
					const [key, value = ''] = pair.split('=').map(decodeURIComponent);
					this._params.push([key, value]);
				}
			}
		} else if (Array.isArray(init)) {
			for (const [key, value] of init) {
				this._params.push([String(key), String(value)]);
			}
		} else if (init && typeof init === 'object') {
			for (const key of Object.keys(init)) {
				this._params.push([key, String(init[key])]);
			}
		}
	}

	append(name, value) { this._params.push([String(name), String(value)]); }
	delete(name) { this._params = this._params.filter(([k]) => k !== name); }
	get(name) { const entry = this._params.find(([k]) => k === name); return entry ? entry[1] : null; }
	getAll(name) { return this._params.filter(([k]) => k === name).map(([, v]) => v); }
	has(name) { return this._params.some(([k]) => k === name); }
	set(name, value) { this.delete(name); this.append(name, value); }
	sort() { this._params.sort(([a], [b]) => a < b ? -1 : a > b ? 1 : 0); }

	*entries() { for (const [k, v] of this._params) yield [k, v]; }
	*keys() { for (const [k] of this._params) yield k; }
	*values() { for (const [, v] of this._params) yield v; }
	[Symbol.iterator]() { return this.entries(); }

	toString() {
		return this._params.map(([k, v]) => encodeURIComponent(k) + '=' + encodeURIComponent(v)).join('&');
	}
}

class URL {
	constructor(url, base) {
		if (base) {
			const baseUrl = new URL(base);
			url = new URL(url, baseUrl.href).href;
		}
		const match = url.match(/^(https?:)\\/\\/([^:/?#]+)(?::(\\d+))?([^?#]*)(\\?[^#]*)?(#.*)?$/);
		if (!match) throw new TypeError('Invalid URL: ' + url);
		this.protocol = match[1];
		this.hostname = match[2];
		this.port = match[3] || '';
		this.pathname = match[4] || '/';
		this.search = match[5] || '';
		this.hash = match[6] || '';
		this.searchParams = new URLSearchParams(this.search);
		this._url = url;
	}

	get origin() { return this.protocol + '//' + this.hostname + (this.port ? ':' + this.port : ''); }
	get host() { return this.hostname + (this.port ? ':' + this.port : ''); }
	get href() {
		let h = this.origin + this.pathname;
		const qs = this.searchParams.toString();
		if (qs) h += '?' + qs;
		if (this.hash) h += this.hash;
		return h;
	}
	toString() { return this.href; }
}

globalThis.URL = URL;
globalThis.URLSearchParams = URLSearchParams;
`
