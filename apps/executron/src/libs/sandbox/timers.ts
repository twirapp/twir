export const timersPolyfill = `
let __nextTimerId = 1;
const __timers = new Map();

globalThis.setTimeout = function(fn, ms) {
	const id = __nextTimerId++;
	__timers.set(id, true);
	const start = Date.now();
	void (async () => {
		while (__timers.has(id)) {
			if (Date.now() - start >= (ms || 0)) {
				__timers.delete(id);
				try { fn(); } catch(e) {}
				return;
			}
			await new Promise(r => r());
		}
	})();
	return id;
};

globalThis.clearTimeout = function(id) {
	__timers.delete(id);
};

globalThis.setInterval = function(fn, ms) {
	const id = __nextTimerId++;
	__timers.set(id, true);
	void (async () => {
		while (__timers.has(id)) {
			await new Promise(r => setTimeout(r, ms || 0));
			if (!__timers.has(id)) return;
			try { fn(); } catch(e) {}
		}
	})();
	return id;
};

globalThis.clearInterval = function(id) {
	__timers.delete(id);
};
`
