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

const messages = ref<ChatMessageType[]>([])

watch(data, (v) => {
	messages.value = v?.chatMessages?.reverse() ?? []
}, {
	immediate: true,
})

const subscription = useChatMessagesSubscription()
watch(subscription.data, (v) => {
	if (v?.chatMessages) {
		messages.value.push(v.chatMessages)
	}
})

onMounted(() => {
	executeQuery({ requestPolicy: 'cache-and-network' })
})

const messagesCount = computed(() => {
	return messages.value.length ?? 0
})

const scrollRef = useTemplateRef<HTMLElement | null>('scrollRef')

const rowVirtualizer = useVirtualizer({
	get count() {
		return messagesCount.value
	},
	getScrollElement: () => scrollRef.value,
	estimateSize: () => 20,
})

const virtualRows = computed(() => rowVirtualizer.value.getVirtualItems())
const totalSize = computed(() => rowVirtualizer.value.getTotalSize())

function measureElement(el: Element) {
	if (!el) {
		return
	}

	rowVirtualizer.value.measureElement(el)

	return undefined
}

watch(messages, (v) => {
	rowVirtualizer.value.scrollToIndex(v.length - 1)
}, {
	immediate: true,
})
</script>

<template>
	<Card>
		<CardContent class="p-0">
			<div
				ref="scrollRef"
				class="overflow-auto h-[75dvh] w-full"
				style="contain: strict"
			>
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
						class="flex flex-col gap-2"
					>
						<div
							v-for="virtualRow in virtualRows"
							:key="virtualRow.key as PropertyKey"
							:ref="measureElement"
							:data-index="virtualRow.index"

							:class="{
								'bg-zinc-800/40': virtualRow.index % 2,
							}"
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
