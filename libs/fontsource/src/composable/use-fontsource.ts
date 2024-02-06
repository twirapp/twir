import { ref, onMounted } from 'vue';

import { loadFontList, loadFont as loadFontById, generateFontKey } from '../api.js';
import type { Font, FontItem } from '../types.js';

export function useFontSource(preloadFonts = true) {
	const loading = ref(true);
	const fontList = ref<FontItem[]>([]);
	const fonts = ref<Font[]>([]);

	onMounted(async () => {
		try {
			if (!preloadFonts) return;
			fontList.value = await loadFontList();
		} catch (err) {
			console.error(err);
		} finally {
			loading.value = false;
		}
	});

	async function loadFont(
		fontId: string,
		fontWeight: number,
		fontStyle: string,
	): Promise<Font | undefined> {
		const fontKey = generateFontKey(fontId, fontWeight, fontStyle);
		for (const fontFace of document.fonts.values()) {
			if (fontFace.family === fontKey) return getFont(fontId);
		}

		try {
			const font = await loadFontById(fontId, fontWeight, fontStyle);
			if (!font) return;
			fonts.value.push(font);
			return font;
		} catch (err) {
			console.error(err);
		}
	}

	function getFont(fontId: string): Font | undefined {
		return fonts.value.find((font) => font.id === fontId);
	}

	return {
		loading,
		fonts,
		fontList,
		loadFont,
		getFont,
		generateFontKey,
	};
}
