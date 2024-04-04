import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useLayout = defineStore('layout', () => {
	const layoutRef = ref<HTMLElement | null>(null);

	function scrollToTop() {
		if (!layoutRef.value) return;
		layoutRef.value.scrollTo({ top: 0, behavior: 'smooth' });
	}

	function scrollToBottom() {
		if (!layoutRef.value) return;
		layoutRef.value.scrollTo({ top: layoutRef.value.scrollHeight, behavior: 'smooth' });
	}

	return {
		layoutRef,
		scrollToTop,
		scrollToBottom,
	};
});
