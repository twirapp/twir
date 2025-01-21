<script setup lang="ts">
import { intlFormat } from 'date-fns'

import type { ChatMessage } from '@/gql/graphql'

const props = defineProps<{
	message: ChatMessage
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
</script>

<template>
	<span>
		<span class="text-xs text-zinc-400">{{ formattedDate }}</span>
		{{ ' ' }}
		<span>
			<span :style="{ color: message.userColor }">{{ message.userDisplayName }}</span>:
		</span>
		{{ ' ' }}
		<span class="break-words">{{ message.text }}</span>
	</span>
</template>
