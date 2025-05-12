<script setup lang="ts">
import { useVirtualizer } from '@tanstack/vue-virtual'
import { computed, onMounted, ref, useTemplateRef, watch } from 'vue'

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

const boxRef = useTemplateRef('boxRef')
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

const totalMessages = computed(() => messages.value.length)

const rowVirtualizer = useVirtualizer({
	get count() {
		return totalMessages.value
	},
	getScrollElement: () => boxRef.value,
	estimateSize: () => 30,
	overscan: 5,
})
const virtualRows = computed(() => rowVirtualizer.value.getVirtualItems())
const totalSize = computed(() => rowVirtualizer.value.getTotalSize())
</script>

<template>
	<Card>
		<CardContent class="p-2 flex flex-col gap-2 h-[80dvh]">
			<div ref="boxRef" class="overflow-y-auto h-full flex-1">
				<div v-if="data?.chatMessages.length === 0">
					No data
				</div>
				<div
					:style="{
						height: `${totalSize}px`,
						width: '100%',
						position: 'relative',
					}"
				>
					<div
						:style="{
							position: 'absolute',
							top: 0,
							left: 0,
							width: '100%',
							transform: `translateY(${virtualRows[0]?.start ?? 0}px)`,
						}"
					>
						<div
							v-for="virtualRow in virtualRows"
							:key="virtualRow.index"
							class="border-b border-border px-2 py-0.5 flex items-center justify-between"
						>
							<ChatMessage
								:message="messages[virtualRow.index]"
							/>
						</div>
					</div>
				</div>
			</div>
		</CardContent>
	</Card>
</template>
