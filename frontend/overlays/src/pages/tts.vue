<script lang="ts" setup>
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import type { TTSSayMessage } from '@/types.js'

import { useTTSOverlayGraphQL } from '@/composables/tts/use-tts-graphql.js'
import { openApi } from '@/api.js'

declare global {
	interface Window {
		webkitAudioContext: typeof AudioContext
	}
}

const isProcessing = ref(false)
const queueMessages = ref<TTSSayMessage[]>([])
const currentAudioBuffer = ref<AudioBufferSourceNode | null>(null)

const route = useRoute()
const apiKey = route.params.apiKey as string

const { connect, destroy } = useTTSOverlayGraphQL({
	onSay: (message) => {
		queueMessages.value.push(message)
		processQueue()
	},
	onSkip: () => {
		currentAudioBuffer.value?.stop()
	},
})

onMounted(() => {
	connect(apiKey)
})

onBeforeUnmount(() => {
	destroy()
})

async function processQueue() {
	if (isProcessing.value) return

	const message = queueMessages.value.shift()
	if (!message) return

	isProcessing.value = true
	await sayMessage(message)
	isProcessing.value = false

	// Process the next item in the queue
	processQueue()
}

async function sayMessage(data: TTSSayMessage) {
	if (!data.text) return

	const audioContext = new (window.AudioContext || window.webkitAudioContext)()
	const gainNode = audioContext.createGain()

	const { data: response } = await openApi.v1.ttsSay({
		voice: data.voice,
		text: data.text,
		volume: Number(data.volume),
		pitch: Number(data.pitch),
		rate: Number(data.rate),
	})

	if (!response) return

	const source = audioContext.createBufferSource()
	currentAudioBuffer.value = source

	// The response is a blob/arraybuffer directly from the API
	const arrayBuffer = await response.arrayBuffer()

	source.buffer = await audioContext.decodeAudioData(arrayBuffer)

	gainNode.gain.value = Number.parseInt(data.volume) / 100
	source.connect(gainNode)
	gainNode.connect(audioContext.destination)

	return new Promise<void>((resolve) => {
		source.onended = () => {
			currentAudioBuffer.value = null
			resolve()
		}

		source.start(0)
	})
}
</script>

<template>
	<!-- TTS overlay has no visible UI, it only plays audio -->
</template>
