<script setup lang="ts">
import { ChatBox } from '@twir/frontend-chat'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRoute } from 'vue-router'

import type { Message } from '@twir/frontend-chat'

import { useChatOverlaySocket } from '@/composables/chat/use-chat-overlay-socket.js'
import {
	type ChatMessage,
	type ChatSettings,
	knownBots,
	useChatTmi,
} from '@/composables/tmi/use-chat-tmi.js'

const route = useRoute()

const messages = ref<Message[]>([])
const maxMessages = ref(30)

const { settings, connect, destroy } = useChatOverlaySocket()

function removeMessageByInternalId(id: string) {
	messages.value = messages.value.filter(m => m.internalId !== id)
}

function removeMessageById(id: string) {
	messages.value = messages.value.filter(m => m.id !== id)
}

function removeMessageByUserName(userName: string) {
	messages.value = messages.value.filter(m => m.sender !== userName)
}

function onMessage(m: ChatMessage) {
	if (m.sender && settings.value.hideBots && knownBots.has(m.sender)) {
		return
	}

	if (settings.value.hideCommands && m.chunks.at(0)?.value.startsWith('!')) {
		return
	}

	const internalId = crypto.randomUUID()

	const showDelay = settings.value.messageShowDelay ?? settings.value.messageShowDelay

	if (messages.value.length >= maxMessages.value) {
		messages.value = messages.value.slice(1)
	}

	setTimeout(() => {
		messages.value.push({
			...m,
			isItalic: m.isItalic ?? false,
			createdAt: new Date(),
			internalId,
			isAnnounce: m.isAnnounce ?? false,
		})
	}, showDelay * 1000)

	const hideTimeout = m.messageHideTimeout ?? settings.value.messageHideTimeout

	if (hideTimeout) {
		setTimeout(() => removeMessageByInternalId(internalId), hideTimeout * 1000)
	}
}

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: settings.value.channelId,
		channelName: settings.value.channelName,
		emotes: {
			ffz: true,
			bttv: true,
			sevenTv: true,
		},
		onMessage,
		onRemoveMessage: removeMessageById,
		onRemoveMessageByUser: removeMessageByUserName,
		onChatClear: () => {
			messages.value = []
		},
	}
})

const chatTmiStore = useChatTmi(chatSettings)

onMounted(() => {
	document.body.style.overflow = 'hidden'

	const apiKey = route.params.apiKey as string
	const overlayId = route.query.id as string
	connect(apiKey, overlayId)
})

onUnmounted(async () => {
	destroy()
	chatTmiStore.destroy()
})
</script>

<template>
	<ChatBox :messages="messages" :settings="settings" />
</template>
