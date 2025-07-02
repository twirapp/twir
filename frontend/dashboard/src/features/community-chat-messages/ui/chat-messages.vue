<script setup lang="ts">
import { useVirtualizer } from '@tanstack/vue-virtual'
import { watchThrottled } from '@vueuse/core'
import { ArrowDownToLine, MoveDown } from 'lucide-vue-next'
import { VNodeRef, computed, nextTick, onMounted, ref, useTemplateRef, watch } from 'vue'

import ChatMessage from './message.vue'
import { useChatMessagesFilters } from '../composables/use-filters'

import {
	type ChatMessage as ChatMessageType,
	useChatMessages,
	useChatMessagesSubscription,
} from '@/api/chat-messages'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'

const filters = useChatMessagesFilters()

const { data, executeQuery } = useChatMessages(filters.computedFilters)

const boxRef = useTemplateRef('boxRef')
const messages = ref<ChatMessageType[]>([])
const isAutoScrolling = ref(true)

watch(
	data,
	(v) => {
		messages.value = v?.chatMessages?.reverse() ?? []
	},
	{
		immediate: true,
	}
)

const subscription = useChatMessagesSubscription()

onMounted(() => {
	executeQuery({ requestPolicy: 'cache-and-network' })
})

const totalMessages = computed(() => messages.value.length)

const rowVirtualizer = useVirtualizer({
	get count() {
		return totalMessages.value
	},
	getScrollElement: () => boxRef.value,
	estimateSize: () => 28,
	overscan: 5,
})

const totalSize = computed(() => rowVirtualizer.value.getTotalSize())
const virtualRows = computed(() => rowVirtualizer.value.getVirtualItems())

function scrollToBottom() {
	rowVirtualizer.value.scrollToIndex(messages.value.length - 1, { align: 'end' })

	if (!isAutoScrolling.value) {
		isAutoScrolling.value = true
	}
}
watch(subscription.data, (v) => {
	if (!v?.chatMessages) return

	messages.value.push(v.chatMessages)
	if (isAutoScrolling.value) {
		scrollToBottom()
	}
})

watch(
	messages,
	() => {
		nextTick(scrollToBottom)
	},
	{ once: true }
)

watchThrottled(
	rowVirtualizer,
	(v) => {
		const currentScrollPosition = (v.scrollRect?.height || 0) + (v.scrollOffset || 0)

		if (v.isScrolling && currentScrollPosition > (v.scrollOffset ?? 0)) {
			isAutoScrolling.value = false
		}

		if (currentScrollPosition >= totalSize.value - 100) {
			isAutoScrolling.value = true
		}
	},
	{ throttle: 500 }
)

function measureElement(el: HTMLDivElement): VNodeRef | undefined {
	if (!el) {
		return
	}

	rowVirtualizer.value.measureElement(el)

	return undefined
}
</script>

<template>
	<Card>
		<CardContent class="p-2 flex flex-col gap-2 h-[78dvh] relative">
			<div v-if="!isAutoScrolling" class="absolute z-10 inset-x-0 bottom-3 w-fit mx-auto">
				<Button variant="secondary" @click="scrollToBottom">
					Scroll to bottom
					<ArrowDownToLine />
				</Button>
			</div>
			<MoveDown v-if="messages.length !== 0" class="absolute top-4 right-4 opacity-40" />
			<div
				ref="boxRef"
				class="overflow-y-auto h-full flex-1 [contain:strict] [overflow-anchor:none]"
			>
				<div v-if="messages.length === 0">No data</div>
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
							:ref="measureElement"
							:data-index="virtualRow.index"
							class="border-b border-border px-2 py-0.5 flex items-center justify-between"
						>
							<ChatMessage :message="messages[virtualRow.index]" />
						</div>
					</div>
				</div>
			</div>
		</CardContent>
	</Card>
</template>
