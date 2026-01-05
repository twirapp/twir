import { onMounted, ref } from 'vue'

import type { Font, FontItem } from '../types'

import { generateFontKey, loadFont as loadFontById, loadFontList } from '../api'

// https://bugzilla.mozilla.org/show_bug.cgi?id=1729089
// https://bugzilla.mozilla.org/show_bug.cgi?id=1780657
// https://sidneyliebrand.io/blog/fixing-font-face-set-entries-not-iterable-in-firefox
function firefoxWorkaroundIterFonts(target: FontFaceSet) {
	const iterable = target.entries()
	const results = []
	let iterator = iterable.next()
	while (iterator.done === false) {
		results.push(iterator.value[0])
		iterator = iterable.next()
	}
	return results
}

export function useFontSource(preloadFonts = true) {
	const loading = ref(true)
	const fontList = ref<FontItem[]>([])
	const fonts = ref<Font[]>([])

	onMounted(async () => {
		try {
			if (!preloadFonts) return
			fontList.value = await loadFontList()
		} catch (err) {
			console.error(err)
		} finally {
			loading.value = false
		}
	})

	async function loadFont(
		fontId: string,
		fontWeight: number,
		fontStyle: string
	): Promise<Font | undefined> {
		const fontKey = generateFontKey(fontId, fontWeight, fontStyle)

		for (const fontFace of firefoxWorkaroundIterFonts(document.fonts)) {
			if (fontFace.family === fontKey) return getFont(fontId)
		}

		try {
			const font = await loadFontById(fontId, fontWeight, fontStyle)
			if (!font) return
			fonts.value.push(font)
			return font
		} catch (err) {
			console.error(err)
		}
	}

	function getFont(fontId: string): Font | undefined {
		return fonts.value.find((font) => font.id === fontId)
	}

	return {
		loading,
		fonts,
		fontList,
		loadFont,
		getFont,
		generateFontKey,
	}
}
