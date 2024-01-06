<script lang="ts" setup>
import { useWebSocket } from '@vueuse/core';
import { ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import { unprotectedApiClient } from '@/api.js';
import { generateSocketUrlWithParams } from '@/helpers.js';

declare global {
	interface Window {
		webkitAudioContext: typeof AudioContext;
	}
}

const queue = ref<Array<Record<string, string>>>([]);
const currentAudioBuffer = ref<AudioBufferSourceNode | null>(null);

const route = useRoute();

const apiKey = route.params.apiKey as string;
const ttsUrl = generateSocketUrlWithParams('/overlays/tts', {
	apiKey,
});

const { data } = useWebSocket(ttsUrl, {
	autoReconnect: {
		delay: 500,
	},
});

watch(data, (message) => {
	const parsedData = JSON.parse(message);
	if (parsedData.eventName === 'say') {
		queue.value.push(parsedData.data);

		if (queue.value.length === 1) {
			processQueue();
		}
	}

	if (parsedData.eventName === 'skip') {
		currentAudioBuffer.value?.stop();
	}
});

const processQueue = async () => {
	if (queue.value.length === 0) {
		return;
	}

	await say(queue.value[0]);
	queue.value = queue.value.slice(1);

	// Process the next item in the queue
	processQueue();
};

const say = async (data: Record<string, string>) => {
	if (!apiKey || !data.text) return;
	const audioContext = new (window.AudioContext || window.webkitAudioContext)();
	const gainNode = audioContext.createGain();

	const req = await unprotectedApiClient.modulesTTSSay({
		voice: data.voice,
		text: data.text,
		volume: Number(data.volume),
		pitch: Number(data.pitch),
		rate: Number(data.rate),
	});

	const source = audioContext.createBufferSource();
	currentAudioBuffer.value = source;

	source.buffer = await audioContext.decodeAudioData(req.response.file.buffer);

	gainNode.gain.value = parseInt(data.volume) / 100;
	source.connect(gainNode);
	gainNode.connect(audioContext.destination);

	return new Promise((resolve) => {
		source.onended = () => {
			currentAudioBuffer.value = null;
			resolve(null);
		};

		source.start(0);
	});
};
</script>
