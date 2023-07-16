
import { useWebSocket } from '@vueuse/core';
import { computed, onMounted, ref, watch } from 'vue';

import { useProfile } from '@/api/index.js';

export type Video = {
	id: string
	channelId: string
	createdAt: string
	deletedAt: string | null
	duration: number
	orderedByDisplayName: string
	orderedById: string
	orderedByName: string
	queuePosition: number
	songLink: null | string
	title: string
	videoId: string
}

export const useYoutubeSocket = () => {
	const videos = ref<Video[]>([]);
	const currentVideo = computed(() => videos.value[0]);

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

	watch(websocket.data, (data) => {
		const parsedData = JSON.parse(data);
			if (parsedData.eventName === 'currentQueue') {
					const incomingVideos = parsedData.data;

					videos.value = incomingVideos;
			}
	});

	onMounted(() => {
		websocket.open();
	});

	const nextVideo = () => {
		// TODO: send socket we deleted video
		videos.value = videos.value.slice(1);
	};

	const deleteVideo = (id: string) => {
		// TODO: send socket we deleted video
		videos.value = videos.value.filter(video => video.id !== id);
	};

	return {
		videos,
		currentVideo,
		nextVideo,
		deleteVideo,
	};
};

