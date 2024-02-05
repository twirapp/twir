import { defineStore } from 'pinia';
import { ref, toRaw, watch } from 'vue';

import { useDudesForm } from './use-dudes-form';

export const useDudesIframe = defineStore('dudes-iframe', () => {
	const { data } = useDudesForm();
	const dudesIframe = ref<HTMLIFrameElement | null>(null);

	function sendIframeMessage(action: string, data?: any) {
		if (!dudesIframe.value) return;
		const payload = JSON.stringify({
			action,
			data: toRaw(data),
		});
		dudesIframe.value.contentWindow?.postMessage(payload);
	}

	watch(dudesIframe, (iframe) => {
		if (!iframe) return;
		iframe.contentWindow?.addEventListener('message', (event) => {
			const parsedData = JSON.parse(event.data);
			if (parsedData.action !== 'get-settings') return;
			sendIframeMessage('settings', data.value);
		});
	});

	return {
		dudesIframe,
		sendIframeMessage,
	};
});
