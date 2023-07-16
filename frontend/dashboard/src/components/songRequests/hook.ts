import { useVideoQueue } from '@mellkam/vue-plyr-queue';
import { useWebSocket } from '@vueuse/core';
import Plyr from 'plyr';
import { computed, onMounted, type ComputedRef } from 'vue';

import { useProfile } from '@/api/index.js';

export const useQueue = (plyr: ComputedRef<Plyr | null>) => {
	const { data: userProfile } = useProfile();

	const socketUrl = computed(() => {
		if (!userProfile?.value) return;

		const host = window.location.host;
		const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
		return `${protocol}://${host}/socket/youtube?apiKey=${userProfile.value.apiKey}`;
	});
	const websocket = useWebSocket(socketUrl, {
		immediate: false,
		autoReconnect: {
			delay: 1000,
		},
	});

	onMounted(() => {
		websocket.open();
	});

	const controls = useVideoQueue({
		plyr,
		initialQueue: [
			{ id: '1', src: 'https://www.youtube.com/watch?v=2-1ymGpV_1A&list=LL&index=4' },
			{ id: '1', src: 'https://www.youtube.com/watch?v=P4ALDytLAXQ' },
		],
		defaultProvider: 'youtube',
	});

	return {
		currentVideo: controls.currentVideo,
		queue: controls.queue,
		nextVideo: controls.nextVideo,
	};
};
