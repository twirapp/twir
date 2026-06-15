<script setup lang="ts">
import Table from '~~/layers/dashboard/app/components/table.vue'
import {
	useModerationWallTable,
} from '~~/layers/dashboard/app/features/moderation/composables/use-moderation-wall-table.js'
import ChatWallSettings from '~~/layers/dashboard/app/features/moderation/ui/chat-wall-settings.vue'

const table = useModerationWallTable()
</script>

<template>
	<div class="flex flex-col gap-4">
		<ChatWallSettings />
		<Table :table="table.table" :is-loading="table.isLoading.value" />
	</div>
</template>
