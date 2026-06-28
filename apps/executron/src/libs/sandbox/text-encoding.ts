export const textEncodingPolyfill = `
class TextEncoder {
	constructor() { this.encoding = 'utf-8'; }

	encode(input = '') {
		const str = String(input);
		const bytes = [];
		for (let i = 0; i < str.length; i++) {
			let code = str.charCodeAt(i);
			if (code < 0x80) {
				bytes.push(code);
			} else if (code < 0x800) {
				bytes.push(0xc0 | (code >> 6), 0x80 | (code & 0x3f));
			} else if (code < 0xd800 || code >= 0xe000) {
				bytes.push(0xe0 | (code >> 12), 0x80 | ((code >> 6) & 0x3f), 0x80 | (code & 0x3f));
			} else {
				i++;
				code = 0x10000 + (((code & 0x3ff) << 10) | (str.charCodeAt(i) & 0x3ff));
				bytes.push(
					0xf0 | (code >> 18),
					0x80 | ((code >> 12) & 0x3f),
					0x80 | ((code >> 6) & 0x3f),
					0x80 | (code & 0x3f)
				);
			}
		}
		return new Uint8Array(bytes);
	}

	encodeInto(src, dest) {
		const encoded = this.encode(src);
		const len = Math.min(encoded.length, dest.length);
		dest.set(encoded.subarray(0, len));
		return { read: src.length, written: len };
	}
}

class TextDecoder {
	constructor(encoding = 'utf-8') { this.encoding = encoding.toLowerCase().replace('-', ''); }

	decode(input) {
		if (!input) return '';
		const bytes = input instanceof ArrayBuffer ? new Uint8Array(input) : input;
		let result = '';
		for (let i = 0; i < bytes.length; i++) {
			const b = bytes[i];
			if (b < 0x80) {
				result += String.fromCharCode(b);
			} else if (b < 0xe0) {
				result += String.fromCharCode(((b & 0x1f) << 6) | (bytes[++i] & 0x3f));
			} else if (b < 0xf0) {
				result += String.fromCharCode(((b & 0x0f) << 12) | ((bytes[++i] & 0x3f) << 6) | (bytes[++i] & 0x3f));
			} else {
				const cp = ((b & 0x07) << 18) | ((bytes[++i] & 0x3f) << 12) | ((bytes[++i] & 0x3f) << 6) | (bytes[++i] & 0x3f);
				result += String.fromCodePoint(cp);
			}
		}
		return result;
	}
}

globalThis.TextEncoder = TextEncoder;
globalThis.TextDecoder = TextDecoder;
`
