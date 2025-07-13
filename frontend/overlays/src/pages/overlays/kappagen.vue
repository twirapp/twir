<script setup lang="ts">
import KappagenOverlay from '@twirapp/kappagen'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import type { Emote, KappagenAnimations, KappagenMethods } from '@twirapp/kappagen/types'

import { useKappagenEmotesBuilder } from '@/composables/kappagen/use-kappagen-builder.js'
import { useKappagenSettings } from '@/composables/kappagen/use-kappagen-settings.js'
import { useKappagenOverlaySocket } from '@/composables/kappagen/use-kappagen-socket.js'
import { type ChatMessage, type ChatSettings, useChatTmi } from '@/composables/tmi/use-chat-tmi.js'

const kappagen = ref<KappagenMethods>()
const route = useRoute()
const emotesBuilder = useKappagenEmotesBuilder()
const { settings, setSettings } = useKappagenSettings()

const socket = useKappagenOverlaySocket({
	playAnimation,
	showEmotes,
	clear: () => kappagen.value?.clear(),
}, emotesBuilder)

watch(socket.settings, (v) => {
	if (!v) return

	setSettings({
		animation: {
			fade: {
				in: v.overlaysKappagen?.animation?.fadeIn ? 1 : 0,
				out: v.overlaysKappagen?.animation?.fadeOut ? 1 : 0,
			},
			zoom: {
				in: v.overlaysKappagen?.animation?.zoomIn ? 1 : 0,
				out: v.overlaysKappagen?.animation?.zoomOut ? 1 : 0,
			},
		},
		cube: {},
		in: {
			fade: v.overlaysKappagen?.animation?.fadeIn,
			zoom: v.overlaysKappagen?.animation?.zoomIn,
		},
		out: {
			fade: v.overlaysKappagen?.animation?.fadeOut,
			zoom: v.overlaysKappagen?.animation?.zoomOut,
		},
		size: {
			max: v.overlaysKappagen?.size?.max ?? 100,
			min: v.overlaysKappagen?.size?.min ?? 50,
		},
		max: v.overlaysKappagen?.emotes.max ?? 100,
		queue: v.overlaysKappagen?.emotes.queue ?? 100,
		time: v.overlaysKappagen?.emotes.time ?? 5000,
		emojiStyle: v.overlaysKappagen.emotes?.emojiStyle,
		excludedEmotes: v.overlaysKappagen?.excludedEmotes ?? [],
	})
})

function playAnimation(emotes: Emote[], animation: KappagenAnimations) {
	if (!kappagen.value) return Promise.resolve()
	return kappagen.value.playAnimation(emotes, animation)
}

function showEmotes(emotes: Emote[]) {
	if (!kappagen.value) return
	kappagen.value.showEmotes(emotes)
}

function onMessage(msg: ChatMessage): void {
	if (msg.type === 'system' || !socket.settings.value?.overlaysKappagen.enableSpawn) return

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
		channelId: socket.settings.value?.overlaysKappagen?.channel.id ?? '',
		channelName: socket.settings.value?.overlaysKappagen?.channel.login ?? '',
		emotes: {
			ffz: socket.settings.value?.overlaysKappagen?.emotes.ffzEnabled,
			bttv: socket.settings.value?.overlaysKappagen?.emotes.bttvEnabled,
			sevenTv: socket.settings.value?.overlaysKappagen?.emotes.sevenTvEnabled,
		},
		onMessage,
	}
})

const { destroy: destroyChat } = useChatTmi(chatSettings)

onMounted(() => {
	const apiKey = route.params.apiKey as string
	if (!apiKey) {
		console.error('API key is required for Kappagen overlay')

		return
	}
	socket.connect(apiKey)
})

onUnmounted(() => {
	socket.destroy()
	destroyChat()
})
</script>

<template>
	<KappagenOverlay ref="kappagen" :config="settings" :is-rave="socket.settings.value?.overlaysKappagen.enableRave" />
</template>
