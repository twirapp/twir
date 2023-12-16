import { ref, onMounted } from 'vue';

import { loadFontItems, loadFont as loadFontById } from '../helpers.js';
import type { FontItem } from '../types.js';

export function useFontSource() {
	const loading = ref(true);
	const fonts = ref<FontItem[]>([]);

	onMounted(async () => {
		try {
			fonts.value = await loadFontItems();
		} catch (err) {
			console.error(err);
		} finally {
			loading.value = false;
		}
	});

	async function loadFont(fontId: string, fontWeight = 400) {
		for (const fontFace of document.fonts.values()) {
			if (fontFace.family === fontId) return;
		}

		try {
			await loadFontById(fontId, fontWeight);
		} catch (err) {
			console.error(err);
		}
	}

	return {
		loading,
		fonts,
		loadFont,
	};
}
