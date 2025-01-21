<script setup lang="ts">
import { useIntervalFn } from '@vueuse/core'

import ChatMessage from './message.vue'
import { useChatMessagesFilters } from '../composables/use-filters'

import { useChatMessages } from '@/api/chat-messages'
import { Card, CardContent } from '@/components/ui/card'

const filters = useChatMessagesFilters()

const { data, executeQuery } = useChatMessages(filters.computedFilters)

useIntervalFn(() => {
	executeQuery({
		requestPolicy: 'network-only',
	})
}, 1000)
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
