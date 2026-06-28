export const base64Polyfill = `
globalThis.btoa = function(str) {
	const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
	let result = '';
	for (let i = 0; i < str.length; i += 3) {
		const a = str.charCodeAt(i);
		const b = i + 1 < str.length ? str.charCodeAt(i + 1) : 0;
		const c = i + 2 < str.length ? str.charCodeAt(i + 2) : 0;
		result += chars[a >> 2];
		result += chars[((a & 3) << 4) | (b >> 4)];
		result += i + 1 < str.length ? chars[((b & 15) << 2) | (c >> 6)] : '=';
		result += i + 2 < str.length ? chars[c & 63] : '=';
	}
	return result;
};

globalThis.atob = function(encoded) {
	const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/';
	encoded = encoded.replace(/[^A-Za-z0-9+/]/g, '');
	let result = '';
	for (let i = 0; i < encoded.length; i += 4) {
		const a = chars.indexOf(encoded[i]);
		const b = chars.indexOf(encoded[i + 1]);
		const c = chars.indexOf(encoded[i + 2]);
		const d = chars.indexOf(encoded[i + 3]);
		result += String.fromCharCode((a << 2) | (b >> 4));
		if (encoded[i + 2] !== '=') result += String.fromCharCode(((b & 15) << 4) | (c >> 2));
		if (encoded[i + 3] !== '=') result += String.fromCharCode(((c & 3) << 6) | d);
	}
	return result;
};
`
