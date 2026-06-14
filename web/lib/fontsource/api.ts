import type { Font, FontItem } from './types'

import { FONTSOURCE_API_URL } from './constants'

export async function loadFontList(): Promise<FontItem[]> {
	const response = await fetch(
		`${FONTSOURCE_API_URL}/fonts?subsets=latin,cyrillic&styles=normal,italic`
	)
	const fonts = (await response.json()) as FontItem[]
	return fonts
}

export async function loadFont(
	fontId: string,
	fontWeight: number,
	fontStyle: string
): Promise<Font | undefined> {
	if (!fontId) return

	const response = await fetch(`${FONTSOURCE_API_URL}/fonts/${fontId}`)
	const font = (await response.json()) as Font

	for (const subset of font.subsets) {
		const fontSource = `url(${font.variants[fontWeight][fontStyle][subset].url.woff2})`
		const fontKey = generateFontKey(fontId, fontWeight, fontStyle)
		const fontFace = new FontFace(fontKey, fontSource)
		await fontFace.load()
		document.fonts.add(fontFace)
	}

	return font
}

export function generateFontKey(fontId: string, weight: number, style: string): string {
	return `${fontId}-${weight}-${style}`
}
