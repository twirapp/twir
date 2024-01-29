import { defineStore } from 'pinia';
import { ref } from 'vue';

export const useDudesIframe = defineStore('dudes-iframe', () => {
	const dudesIframe = ref<HTMLIFrameElement | null>(null);
	return { dudesIframe };
});
