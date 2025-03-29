<script setup lang="ts">
import type { ChatWall } from '@/api/moderation-chat-wall.ts'

defineProps<{
	chatWall: ChatWall
}>()
</script>

<template>
	<div>
		<span>
			{{ chatWall.action.at(0) + chatWall.action.slice(1).toLowerCase() }}
		</span>
		<span v-if="chatWall.action === 'TIMEOUT' && chatWall.timeoutDurationSeconds">
			({{ chatWall.timeoutDurationSeconds }}sec)
		</span>
	</div>
</template>
