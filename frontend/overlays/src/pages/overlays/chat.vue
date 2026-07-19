<script setup lang="ts">
import type { Message, MessagePlatform } from '@twir/frontend-chat'

import { ChatBox } from '@twir/frontend-chat'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

import { useFragmentsToChunks } from '@/composables/chat/fragments-to-chunks.js'
import { useChatOverlaySocket } from '@/composables/chat/use-chat-overlay-socket.js'
import { knownBots } from '@/composables/tmi/use-chat-tmi.js'
import { useThirdPartyEmotes } from '@/composables/tmi/use-third-party-emotes.js'
import { ChatOverlayModerationEventType } from '@/gql/graphql.js'

type IncomingMessage = Omit<Message, 'createdAt' | 'internalId'>

const messages = ref<Message[]>([])
const maxMessages = ref(30)

const {
	overlaySettings: settings,
	neededData,
	chatLibSettings,
	chatMessages,
	chatModerationEvents,
} = useChatOverlaySocket()
const { fragmentsToChunks } = useFragmentsToChunks()

const thirdPartyEmotesSettings = computed(() => ({
	channelId: neededData.value?.authenticatedUser.twitchProfile?.id ?? '',
	emotes: {
		ffz: true,
		bttv: true,
		sevenTv: true,
	},
}))

const { destroy: destroyThirdPartyEmotes } = useThirdPartyEmotes(thirdPartyEmotesSettings)

function removeMessageByInternalId(id: string) {
	messages.value = messages.value.filter((m) => m.internalId !== id)
}

function removeMessageById(id: string) {
	messages.value = messages.value.filter((m) => m.id !== id)
}

function removeMessageByUserName(userName: string) {
	const normalized = userName.toLowerCase()
	messages.value = messages.value.filter((m) => m.sender?.toLowerCase() !== normalized)
}

function pushMessage(m: IncomingMessage) {
	if (m.sender && settings.value?.hideBots && knownBots.has(m.sender.toLowerCase())) {
		return
	}

	if (settings.value?.hideCommands && m.chunks.at(0)?.value.startsWith('!')) {
		return
	}

	const internalId = crypto.randomUUID()

	const showDelay = settings.value?.messageShowDelay || 0

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

	const hideTimeout = settings.value?.messageHideTimeout

	if (hideTimeout) {
		setTimeout(() => removeMessageByInternalId(internalId), hideTimeout * 1000)
	}
}

watch(chatMessages, (v) => {
	const event = v?.chatMessagesByApiKey
	if (!event) return

	const platform = event.platform as MessagePlatform
	const isKick = platform === 'kick'

	const badges: Record<string, string> = {}
	const kickBadges: Array<{ type: string; text: string }> = []
	for (const badge of event.badges) {
		if (isKick) {
			kickBadges.push({ type: badge.setId, text: badge.text ?? badge.setId })
		} else if (badge.versionId) {
			badges[badge.setId] = badge.versionId
		}
	}

	pushMessage({
		id: event.messageId ?? event.id,
		type: 'message',
		platform,
		chunks: fragmentsToChunks(event.fragments),
		sender: event.userName,
		senderColor: event.userColor || undefined,
		senderDisplayName: event.userDisplayName,
		badges: isKick ? undefined : badges,
		kickBadges: isKick ? kickBadges : undefined,
		isItalic: false,
		isAnnounce: event.messageType === 'announcement',
		announceColor: event.announceColor ?? undefined,
	})
})

watch(chatModerationEvents, (v) => {
	const event = v?.overlaysChatModerationEvents
	if (!event) return

	switch (event.type) {
		case ChatOverlayModerationEventType.UserBanned:
			if (event.userLogin) {
				removeMessageByUserName(event.userLogin)
			}
			break
		case ChatOverlayModerationEventType.MessageDeleted:
			if (event.deletedMessageId) {
				removeMessageById(event.deletedMessageId)
			}
			break
		case ChatOverlayModerationEventType.ChatCleared:
			messages.value = []
			break
	}
})

onMounted(() => {
	document.body.style.overflow = 'hidden'
})

onUnmounted(() => {
	destroyThirdPartyEmotes()
})
</script>

<template>
	<ChatBox
		v-if="chatLibSettings"
		:messages
		:settings="chatLibSettings"
	/>
</template>
