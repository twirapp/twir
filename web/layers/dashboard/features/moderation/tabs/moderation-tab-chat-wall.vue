<script setup lang="ts">
import Table from '#layers/dashboard/components/table.vue'
import {
	useModerationWallTable,
} from '~/features/moderation/composables/use-moderation-wall-table.ts'
import ChatWallSettings from '~/features/moderation/ui/chat-wall-settings.vue'

const table = useModerationWallTable()
</script>

<template>
	<div class="flex flex-col gap-4">
		<ChatWallSettings />
		<Table :table="table.table" :is-loading="table.isLoading.value" />
	</div>
</template>
