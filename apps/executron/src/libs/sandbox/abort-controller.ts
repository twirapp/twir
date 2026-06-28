export const abortControllerPolyfill = `
class AbortSignal {
	constructor() {
		this.aborted = false;
		this.onabort = null;
		this._listeners = [];
	}

	addEventListener(type, listener) {
		if (type === 'abort') this._listeners.push(listener);
	}

	removeEventListener(type, listener) {
		if (type === 'abort') {
			this._listeners = this._listeners.filter(l => l !== listener);
		}
	}

	_dispatchAbort() {
		this.aborted = true;
		const event = { type: 'abort' };
		if (this.onabort) this.onabort(event);
		for (const listener of this._listeners) listener(event);
	}

	throwIfAborted() {
		if (this.aborted) throw new DOMException('The operation was aborted.', 'AbortError');
	}
}

class AbortController {
	constructor() {
		this.signal = new AbortSignal();
	}

	abort() {
		this.signal._dispatchAbort();
	}
}

globalThis.AbortController = AbortController;
globalThis.AbortSignal = AbortSignal;
`
