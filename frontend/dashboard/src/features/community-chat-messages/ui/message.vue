<script setup lang="ts">
import { intlFormat } from 'date-fns'
import { computed } from 'vue'

import type { ChatMessage } from '@/gql/graphql'

const props = defineProps<{
	message: ChatMessage
	withChannel?: boolean
}>()

const createdAt = new Date(props.message.createdAt)
const formattedDate = intlFormat(createdAt, {
	localeMatcher: 'best fit',
	hour: 'numeric',
	minute: 'numeric',
	second: 'numeric',
	day: 'numeric',
	month: 'numeric',
	year: 'numeric',
})

const parsedMessage = computed(() => {
	const urlRegex = /(https?:\/\/[^\s]+)/g
	const parts = props.message.text.split(urlRegex)

	return parts.map(part => {
		if (urlRegex.test(part)) {
			return {
				type: 'url',
				content: part,
			}
		}
		return {
			type: 'text',
			content: part,
		}
	})
})
</script>

<template>
	<span>
		<span class="text-xs text-zinc-400">{{ formattedDate }}</span>
		{{ ' ' }}
		<template v-if="withChannel">
			<span class="text-zinc-400">#{{ message.channelLogin }}</span>
			{{ ' | ' }}
		</template>
		<span>
			<span :style="{ color: message.userColor }">{{ message.userDisplayName }}</span>:
		</span>
		{{ ' ' }}
		<span class="wrap-break-word">
			<template v-for="(part, index) in parsedMessage" :key="index">
				<a
					v-if="part.type === 'url'"
					:href="part.content"
					target="_blank"
					rel="noopener noreferrer"
					class="underline"
				>{{ part.content }}</a>
				<span v-else>{{ part.content }}</span>
			</template>
		</span>
	</span>
</template>
