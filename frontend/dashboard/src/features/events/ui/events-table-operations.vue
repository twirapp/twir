<script setup lang="ts">
import { flatOperations } from '@twir/dashboard/src/features/events/constants/helpers.ts'

import type { EventOperation } from '@/api/events.ts'

defineProps<{
	operations: EventOperation[]
}>()
</script>

<template>
	<div class="flex flex-row flex-wrap gap-2">
		<div
			v-for="operation of operations"
			:key="operation.id"
			class="px-2 py-1 rounded-sm text-xs"
			:class="{
				'bg-red-900 text-red-300': flatOperations[operation.type]?.color === 'error',
				'bg-green-900 text-green-300': flatOperations[operation.type]?.color === 'success',
				'bg-blue-900 text-blue-300': flatOperations[operation.type]?.color === 'info',
				'bg-yellow-900 text-yellow-300': flatOperations[operation.type]?.color === 'warning',
			}"
		>
			{{ flatOperations[operation.type].name }}
		</div>
	</div>
</template>
