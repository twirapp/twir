<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'

import ChatMessage from './message.vue'
import { useChatMessagesFilters } from '../composables/use-filters'

import {
	type ChatMessage as ChatMessageType,
	useChatMessages,
	useChatMessagesSubscription,
} from '@/api/chat-messages'
import { Card, CardContent } from '@/components/ui/card'

const filters = useChatMessagesFilters()

const { data, executeQuery } = useChatMessages(filters.computedFilters)

const messages = ref<ChatMessageType[]>([])

watch(data, (v) => {
	messages.value = v?.chatMessages ?? []
}, {
	immediate: true,
})

const subscription = useChatMessagesSubscription()

watch(subscription.data, (v) => {
	if (v?.chatMessages) {
		messages.value.unshift(v.chatMessages)
	}
})

onMounted(() => {
	executeQuery({ requestPolicy: 'cache-and-network' })
})
</script>

<template>
	<Card>
		<CardContent class="p-2 flex flex-col gap-2 overflow-y-auto">
			<div v-if="data?.chatMessages.length === 0">
				No data
			</div>
			<ChatMessage
				v-for="message of data?.chatMessages"
				:key="message.id"
				:message="message"
			/>
		</CardContent>
	</Card>
</template>
