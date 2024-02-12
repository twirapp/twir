import { defineStore } from 'pinia';

interface NowPlayingPostMessage {
	action: string;
	data?: any;
}

export const useNowPlayingIframe = defineStore('now-playing-iframe', () => {
	const isIframe = Boolean(window.frameElement);

	async function onPostMessage(msg: MessageEvent<string>) {
		const parsedData = JSON.parse(msg.data) as NowPlayingPostMessage;
		console.log(parsedData);
	}

	function connect() {
		console.log('attaching');
		if (!isIframe) return;
		window.addEventListener('message', onPostMessage);
	}

	function destroy() {
		if (!isIframe) return;
		window.removeEventListener('message', onPostMessage);
	}

	return {
		connect,
		destroy,
		isIframe,
	};
});
