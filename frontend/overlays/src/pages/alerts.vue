<script lang="ts" setup>
import { useWebSocket } from '@vueuse/core';
import { ref, watch } from 'vue';
import { useRoute } from 'vue-router';

declare global {
  interface Window {
    webkitAudioContext: typeof AudioContext
  }
}

const queue = ref<Array<{
	id: string,
	channel_id: string,
	audio_id: string,
	audio_volume: number
}>>([]);
const currentAudioBuffer = ref<AudioBufferSourceNode | null>(null);

const route = useRoute();
const apiKey = route.params.apiKey as string;

const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
const host = window.location.host;
const socket = useWebSocket(`${protocol}://${host}/socket/alerts?apiKey=${apiKey}`, {
	immediate: true,
	autoReconnect: {
		delay: 500,
	},
});

watch(socket.data, (message) => {
	const parsedData = JSON.parse(message);

	if (parsedData.eventName === 'trigger') {
		queue.value.push(parsedData.data);

		if (queue.value.length === 1) {
			processQueue();
		}
	}
});

const processQueue = async () => {
	if (queue.value.length === 0) {
			return;
		}

		const current = queue.value[0];
		if (current.audio_id) {
			await playAudio(current.channel_id, current.audio_id, current.audio_volume);
		}

		// change next val
		queue.value = queue.value.slice(1);

		// Process the next item in the queue
		processQueue();
};

const playAudio = async (channelId: string, audioId: string, volume: number) => {
	const query = new URLSearchParams({
		channel_id: channelId,
		file_id: audioId,
	});
	const req = await fetch(`${window.location.origin}/api/files/?${query}`);
	if (!req.ok) {
		console.error(await req.text());
		return;
	}

	const audioContext = new (window.AudioContext || window.webkitAudioContext)();
	const gainNode = audioContext.createGain();

	const data = await req.arrayBuffer();

	const source = audioContext.createBufferSource();
	currentAudioBuffer.value = source;

	source.buffer = await audioContext.decodeAudioData(data);

	gainNode.gain.value = volume / 100;
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
