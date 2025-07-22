<script setup lang="ts">
import { flatOperations } from '@/features/events/constants/helpers.ts'

import type { EventOperation } from '@/api/events.ts'
import { getOperationColor } from '@/features/events/composables/use-operation-color.ts'
import { EventOperationType } from '@/gql/graphql.ts'

defineProps<{
	operations: EventOperation[]
}>()

function getOperationTextColor(operationType: EventOperationType): string {
	const operation = flatOperations[operationType]

	switch (operation.color) {
		case 'info':
			return 'text-white'
		case 'success':
		case 'warning':
			return 'text-[#333333]'
		case 'error':
			return 'text-white'
		case 'default':
			return 'text-white'
		default:
			return 'text-white'
	}
}
</script>

<template>
	<div class="flex flex-row flex-wrap gap-2">
		<div
			v-for="operation of operations"
			:key="operation.id"
			class="px-2 py-1 rounded-sm text-xs"
			:class="[getOperationColor(operation.type), getOperationTextColor(operation.type)]"
		>
			{{ flatOperations[operation.type].name }}
		</div>
	</div>
</template>
