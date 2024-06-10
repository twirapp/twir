<script setup lang="ts">
import { ChatBox } from '@twir/frontend-chat'
import { computed, onMounted, onUnmounted, ref } from 'vue'

import type { Message } from '@twir/frontend-chat'

import { useChatOverlaySocket } from '@/composables/chat/use-chat-overlay-socket.js'
import {
	type ChatMessage,
	type ChatSettings,
	knownBots,
	useChatTmi,
} from '@/composables/tmi/use-chat-tmi.js'

const messages = ref<Message[]>([])
const maxMessages = ref(30)

const { overlaySettings: settings, neededData, chatLibSettings } = useChatOverlaySocket()

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
	if (m.sender && settings.value?.hideBots && knownBots.has(m.sender)) {
		return
	}

	if (settings.value?.hideCommands && m.chunks.at(0)?.value.startsWith('!')) {
		return
	}

	const internalId = crypto.randomUUID()

	const showDelay = (settings.value?.messageShowDelay ?? settings.value?.messageShowDelay) || 0

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

	const hideTimeout = m.messageHideTimeout ?? settings.value?.messageHideTimeout

	if (hideTimeout) {
		setTimeout(() => removeMessageByInternalId(internalId), hideTimeout * 1000)
	}
}

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: neededData.value?.authenticatedUser.id ?? '',
		channelName: neededData.value?.authenticatedUser.twitchProfile.login ?? '',
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

	chatTmiStore.destroy()
})

onUnmounted(async () => {
	await chatTmiStore.destroy()
})
</script>

<template>
	<ChatBox
		v-if="chatLibSettings"
		:messages
		:settings="chatLibSettings"
	/>
</template>
