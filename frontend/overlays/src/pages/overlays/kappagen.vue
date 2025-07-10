<script setup lang="ts">
import KappagenOverlay from '@twirapp/kappagen'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import type { Emote, KappagenAnimations, KappagenMethods } from '@twirapp/kappagen/types'

import { useKappagenEmotesBuilder } from '@/composables/kappagen/use-kappagen-builder.js'
import { useKappagenIframe } from '@/composables/kappagen/use-kappagen-iframe.js'
import { useKappagenSettings } from '@/composables/kappagen/use-kappagen-settings.js'
import { useKappagenOverlaySocket } from '@/composables/kappagen/use-kappagen-socket.js'
import { type ChatMessage, type ChatSettings, useChatTmi } from '@/composables/tmi/use-chat-tmi.js'

const {
	connect: connectSocket,
	settings: socketSettings,
	destroy: destroySocket,
} = useKappagenOverlaySocket()
const iframe = useKappagenIframe({
	playAnimation,
	showEmotes,
	clear: () => {
		kappagen.value?.clear()
	},
})

const kappagen = ref<KappagenMethods>()
const route = useRoute()
const { settings, setSettings } = useKappagenSettings()

watch(socketSettings, (v) => {
	if (!v || window.frameElement) return

	setSettings(v)
})

function playAnimation(emotes: Emote[], animation: KappagenAnimations) {
	if (!kappagen.value) return Promise.resolve()
	console.log(emotes, animation)
	return kappagen.value.playAnimation(emotes, animation)
}

function showEmotes(emotes: Emote[]) {
	if (!kappagen.value) return
	kappagen.value.showEmotes(emotes)
}

const emotesBuilder = useKappagenEmotesBuilder()

const socket = useKappagenOverlaySocket({
	playAnimation,
	showEmotes,
	emotesBuilder,
})

function onMessage(msg: ChatMessage): void {
	if (msg.type === 'system' || !overlaySettings.value?.enableSpawn) return

	const firstChunk = msg.chunks.at(0)!
	if (firstChunk.type === 'text' && firstChunk.value.startsWith('!')) {
		return
	}

	const generatedEmotes = emotesBuilder.buildSpawnEmotes(msg.chunks)
	if (!generatedEmotes.length) return
	showEmotes(generatedEmotes)
}

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: settings.value?.channelId ?? '', // todo: take from config
		channelName: settings.value?.channelName ?? '', // todo: take from config
		emotes: {
			ffz: true, // todo: take from config
			bttv: true, // todo: take from config
			sevenTv: true, // todo: take from config
		},
		onMessage,
	}
})

const { destroy: destroyChat } = useChatTmi(chatSettings)

onMounted(() => {
	if (window.frameElement) {
		iframe.create()
	} else {
		const apiKey = route.params.apiKey as string
		connectSocket(apiKey)
	}
})

onUnmounted(() => {
	iframe.destroy()
	destroySocket()
	destroyChat()
})
</script>

<template>
	<KappagenOverlay ref="kappagen" :config="settings" :is-rave="settings?.enableRave" />
</template>
