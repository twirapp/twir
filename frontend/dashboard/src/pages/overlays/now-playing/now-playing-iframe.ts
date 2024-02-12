import { defineStore } from 'pinia';
import { ref, toRaw, watch } from 'vue';

export const useNowPlayingIframe = defineStore('now-playing-iframe', () => {
	const nowPlayingIframe = ref<HTMLIFrameElement | null>(null);

	function sendIframeMessage(action: string, data?: any) {
		if (!nowPlayingIframe.value) return;
		const payload = JSON.stringify({
			action,
			data: toRaw(data),
		});
		nowPlayingIframe.value.contentWindow?.postMessage(payload);
	}

	watch(nowPlayingIframe, (iframe) => {
		if (!iframe) return;

		sendIframeMessage('track', {
			qwe: 'qwe',
		});
	}, { immediate: true });

	return {
		nowPlayingIframe,
		sendIframeMessage,
	};
});
