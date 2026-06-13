import { ref } from 'vue'

export function useDudesIframe() {
	const dudesIframe = ref(null)

	function sendIframeMessage(action: string) {
		// stub
	}

	return { dudesIframe, sendIframeMessage }
}
