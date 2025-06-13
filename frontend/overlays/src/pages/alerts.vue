<script lang="ts" setup>
import { useWebSocket } from '@vueuse/core'
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import { openApi } from '@/api.ts'
import { generateSocketUrlWithParams } from '@/helpers.js'

declare global {
	interface Window {
		webkitAudioContext: typeof AudioContext
	}
}

const queue = ref<Array<{
	id: string
	channel_id: string
	audio_id: string
	audio_volume: number
}>>([])

const currentAudioBuffer = ref<AudioBufferSourceNode | null>(null)
const route = useRoute()

const apiKey = route.params.apiKey as string
const alertsUrl = generateSocketUrlWithParams('/overlays/alerts', {
	apiKey,
})

const socket = useWebSocket(alertsUrl, {
	immediate: true,
	autoReconnect: {
		delay: 500,
	},
})

watch(socket.data, (message) => {
	const parsedData = JSON.parse(message)

	if (parsedData.eventName === 'trigger') {
		queue.value.push(parsedData.data)

		if (queue.value.length === 1) {
			processQueue()
		}
	}
})

const cachedFiles: Record<string, AudioBuffer> = {}

async function processQueue(): Promise<void> {
	if (queue.value.length === 0) {
		return
	}

	const current = queue.value[0]
	if (current.audio_id) {
		await playAudio(current.channel_id, current.audio_id, current.audio_volume)
	}

	// change next val
	queue.value = queue.value.slice(1)

	// Process the next item in the queue
	processQueue()
}

const audioContext = new (window.AudioContext || window.webkitAudioContext)()

async function playAudio(channelId: string, audioId: string, volume: number): Promise<unknown> {
	let data: AudioBuffer
	if (cachedFiles[audioId]) {
		data = cachedFiles[audioId]
	} else {
		const req = await openApi.v1.channelsFilesContentDetail(channelId, audioId)
		if (!req.ok) {
			console.error(await req.text())
			return
		}

		const arrayBuffer = await req.arrayBuffer()

		data = await audioContext.decodeAudioData(arrayBuffer)
		cachedFiles[audioId] = data
	}

	if (!data) {
		console.error('Cannot play audio, no data')
		return
	}

	const gainNode = audioContext.createGain()

	const source = audioContext.createBufferSource()
	currentAudioBuffer.value = source

	source.buffer = data

	gainNode.gain.value = volume / 100
	source.connect(gainNode)
	gainNode.connect(audioContext.destination)

	return new Promise((resolve) => {
		source.onended = () => {
			currentAudioBuffer.value = null
			resolve(null)
		}

		source.start(0)
	})
}
</script>
