<script setup lang="ts">
import type { ChatWall } from '#layers/dashboard/api/moderation-chat-wall.ts'


const { t } = useI18n()

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
			({{ chatWall.timeoutDurationSeconds }}{{ t('chatWall.table.timeout') }})
		</span>
	</div>
</template>
