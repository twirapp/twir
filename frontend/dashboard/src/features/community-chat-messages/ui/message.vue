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

const normalizedPlatform = computed(() => props.message.platform.toLowerCase())
</script>

<template>
	<span>
		<span class="text-xs text-zinc-400">{{ formattedDate }}</span>
		{{ ' ' }}
		<template v-if="withChannel">
			<span class="text-zinc-400">#{{ message.channelLogin }}</span>
			{{ ' | ' }}
		</template>

		<span class="inline-flex items-center gap-1 align-middle mr-1">
			<svg
				v-if="normalizedPlatform === 'kick'"
				xmlns="http://www.w3.org/2000/svg"
				fill="currentColor"
				viewBox="0 0 24 24"
				class="size-3 text-[#53FC18]"
			>
				<path d="M3 5h3.5l5 6.5-5 6.5H3V5z" />
				<path d="M15 5h3v13h-3V5z" />
			</svg>
			<svg
				v-else
				xmlns="http://www.w3.org/2000/svg"
				fill="currentColor"
				viewBox="0 0 24 24"
				class="size-3 text-[#9146FF]"
			>
				<path
					d="M1.3 4.6 2.8.8h19.9v14.5L17 21h-4.6l-3 2.9H6V21H1.3V4.6Zm15.8 12.6 3.7-3.8V2.7H5v14.5h3.7v3l2.8-3h5.6Z"
				/>
				<path d="M17.1 7h-1.8v5.5H17V7Zm-4.6 0h-1.9v5.5h1.9V7Z" />
			</svg>
		</span>

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
