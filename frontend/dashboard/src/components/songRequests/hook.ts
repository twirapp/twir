import { useWebSocket } from '@vueuse/core';
import { defineStore } from 'pinia';
import { computed, onMounted, ref, watch } from 'vue';

import { useProfile, useYoutubeModuleSettings } from '@/api/index.js';

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

export const useYoutubeSocket = defineStore('youtubeSocket', () => {
	const youtubeModuleManager = useYoutubeModuleSettings();
	const { data: youtubeSettings } = youtubeModuleManager.getAll();
	const youtubeModuleUpdater = youtubeModuleManager.update();

	const videos = ref<Video[]>([]);
	const currentVideo = computed(() => videos.value[0]);

	const { data: userProfile } = useProfile();

	async function banUser(userId: string) {
		await youtubeModuleUpdater.mutateAsync({
			data: {
				...youtubeSettings.value!.data!,
				denyList: {
					...youtubeSettings.value!.data!.denyList!,
					users: [
						...youtubeSettings.value!.data!.denyList!.users,
						userId,
					],
				},
			},
		});

		callWsSkip(videos.value.filter(video => video.orderedById === userId).map(video => video.id));
		videos.value = videos.value.filter(video => video.orderedById !== userId);
	}

	async function banSong(videoId: string) {
		await youtubeModuleUpdater.mutateAsync({
			data: {
				...youtubeSettings.value!.data!,
				denyList: {
					...youtubeSettings.value!.data!.denyList!,
					songs: [
						...youtubeSettings.value!.data!.denyList!.songs,
						videoId,
					],
				},
			},
		});

		const video = videos.value.find(video => video.videoId === videoId);

		callWsSkip(video!.id);
		videos.value = videos.value.filter(video => video.videoId !== videoId);
	}

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

		if (parsedData.eventName === 'newTrack') {
			videos.value.push(parsedData.data);
		}

		if (parsedData.eventName === 'removeTrack') {
			videos.value = videos.value.filter(video => video.id !== parsedData.data.id);
		}
	});

	watch(socketUrl, (v) => {
		if (!v) return;
		websocket.open();
	});

	onMounted(() => {
		if (websocket.status.value !== 'OPEN' && socketUrl.value) {
			websocket.open();
		}
	});

	const callWsSkip = (ids: string | string[]) => {
		const request = JSON.stringify({
			eventName: 'skip',
			data: Array.isArray(ids) ? ids : [ids],
		});

		websocket.send(request);
	};

	const nextVideo = () => {
		callWsSkip(currentVideo.value.id);
		videos.value = videos.value.slice(1);
	};

	const deleteVideo = (id: string) => {
		callWsSkip(id);
		videos.value = videos.value.filter(video => video.id !== id);
	};

	const deleteAllVideos = () => {
		callWsSkip(videos.value.map(video => video.id));
		videos.value = [];
	};

	const moveVideo = (id: string, newPosition: number) => {
		const currentIndex = videos.value.findIndex(video => video.id === id);
		const itemToMove = videos.value.splice(currentIndex, 1)[0];
		videos.value.splice(newPosition, 0, itemToMove);

		videos.value.forEach((video, index) => {
			video.queuePosition = index + 1;
		});

		const request = JSON.stringify({
			eventName: 'reorder',
			data: videos.value,
		});
		websocket.send(request);
	};

	const sendPlaying = () => {
		if (!currentVideo.value) return;

		const request = JSON.stringify({
			eventName: 'play',
			data: {
				id: currentVideo.value.id,
				duration: currentVideo.value.duration,
			},
		});
		websocket.send(request);
	};

	return {
		videos,
		currentVideo,
		nextVideo,
		deleteVideo,
		deleteAllVideos,
		moveVideo,
		sendPlaying,
		banUser,
		banSong,
	};
});

