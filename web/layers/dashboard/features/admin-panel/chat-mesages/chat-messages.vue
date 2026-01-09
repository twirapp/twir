<script setup lang="ts">
import { ref, watch } from 'vue'

import type { ChatMessage as ChatMessageType } from '#layers/dashboard/api/chat-messages.ts'

import { useAllChatMessagesSubscription } from '#layers/dashboard/api/admin/chat-messages.ts'

import ChatMessage from '~/features/community-chat-messages/ui/message.vue'

const subscription = useAllChatMessagesSubscription()

const messages = ref<ChatMessageType[]>([])

watch(subscription.data, (v) => {
	if (v?.adminChatMessages) {
		messages.value.unshift(v.adminChatMessages)
	}
})
</script>

<template>
	<UiCard>
		<UiCardContent class="p-2 flex flex-col gap-2 overflow-y-auto min-h-[50dvh]">
			<div v-if="messages.length === 0" class="flex justify-center items-center">
				No data
			</div>
			<ChatMessage
				v-for="message of messages"
				:key="message.id"
				with-channel
				:message="message"
			/>
		</UiCardContent>
	</UiCard>
</template>
