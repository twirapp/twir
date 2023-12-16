import { FONTSOURCE_API_URL } from './constants.js';
import type { Font, FontItem } from './types.js';

export async function loadFontItems() {
	const response = await fetch(`${FONTSOURCE_API_URL}/fonts`);
	const fonts = await response.json() as FontItem[];
	return fonts;
}

export async function loadFont(fontId: string, fontWeight = 400) {
	const response = await fetch(`${FONTSOURCE_API_URL}/fonts/${fontId}`);
	const font = await response.json() as Font;

	// eslint-disable-next-line @typescript-eslint/ban-ts-comment
	// @ts-ignore
	const fontSource = `url(${font.variants[fontWeight].normal.latin.url.woff2})`;
	const fontFace = new FontFace(font.id, fontSource);
	await fontFace.load();
	document.fonts.add(fontFace);
}
