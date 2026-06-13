import { createGlobalState, useWebSocket } from '@vueuse/core'
import { computed, onMounted, ref, watch } from 'vue'

import { useProfile } from '@/api/auth'
import { useSongRequestsApi } from '@/api/song-requests'

export interface Video {
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

export const useYoutubeSocket = createGlobalState(() => {
	const youtubeModuleManager = useSongRequestsApi()
	const { data: youtubeSettings } = youtubeModuleManager.useSongRequestQuery()
	const youtubeModuleUpdater = youtubeModuleManager.useSongRequestMutation()

	const videos = ref<Video[]>([])
	const currentVideo = computed(() => videos.value[0])

	const { data: userProfile } = useProfile()

	const socketUrl = computed(() => {
		if (!userProfile?.value) return

		const host = window.location.host
		const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
		return `${protocol}://${host}/socket/youtube?apiKey=${userProfile.value.apiKey}`
	})

	const websocket = useWebSocket(socketUrl, {
		immediate: false,
		autoReconnect: {
			delay: 1000,
		},
	})

	function callWsSkip(ids: string | string[]) {
		const request = JSON.stringify({
			eventName: 'skip',
			data: Array.isArray(ids) ? ids : [ids],
		})

		websocket.send(request)
	}

	async function banUser(userId: string) {
		if (!youtubeSettings.value?.songRequests) return
		const settings = youtubeSettings.value.songRequests

		await youtubeModuleUpdater.executeMutation({
			opts: {
				...settings,
				denyList: {
					...settings.denyList,
					users: [...settings.denyList.users, userId],
				},
			},
		})

		callWsSkip(
			videos.value.filter((video) => video.orderedById === userId).map((video) => video.id)
		)
		videos.value = videos.value.filter((video) => video.orderedById !== userId)
	}

	async function banSong(videoId: string) {
		if (!youtubeSettings.value?.songRequests) return
		const settings = youtubeSettings.value.songRequests

		await youtubeModuleUpdater.executeMutation({
			opts: {
				...settings,
				denyList: {
					...settings.denyList,
					songs: [...settings.denyList.songs, videoId],
				},
			},
		})

		const video = videos.value.find((video) => video.videoId === videoId)

		callWsSkip(video!.id)
		videos.value = videos.value.filter((video) => video.videoId !== videoId)
	}

	watch(websocket.data, (data) => {
		const parsedData = JSON.parse(data)
		if (parsedData.eventName === 'currentQueue') {
			const incomingVideos = parsedData.data

			videos.value = incomingVideos
		}

		if (parsedData.eventName === 'newTrack') {
			videos.value.push(parsedData.data)
		}

		if (parsedData.eventName === 'removeTrack') {
			videos.value = videos.value.filter((video) => video.id !== parsedData.data.id)
		}
	})

	watch(socketUrl, (v) => {
		if (!v) return
		websocket.open()
	})

	onMounted(() => {
		if (websocket.status.value !== 'OPEN' && socketUrl.value) {
			websocket.open()
		}
	})

	function nextVideo() {
		callWsSkip(currentVideo.value.id)
		videos.value = videos.value.slice(1)
	}

	function deleteVideo(id: string) {
		callWsSkip(id)
		videos.value = videos.value.filter((video) => video.id !== id)
	}

	function deleteAllVideos() {
		callWsSkip(videos.value.map((video) => video.id))
		videos.value = []
	}

	function moveVideo(id: string, newPosition: number) {
		const currentIndex = videos.value.findIndex((video) => video.id === id)
		const itemToMove = videos.value.splice(currentIndex, 1)[0]
		videos.value.splice(newPosition, 0, itemToMove)

		videos.value.forEach((video, index) => {
			video.queuePosition = index + 1
		})

		const request = JSON.stringify({
			eventName: 'reorder',
			data: videos.value,
		})
		websocket.send(request)
	}

	function sendPlaying() {
		if (!currentVideo.value) return

		const request = JSON.stringify({
			eventName: 'play',
			data: {
				id: currentVideo.value.id,
				duration: currentVideo.value.duration,
			},
		})
		websocket.send(request)
	}

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
	}
})
