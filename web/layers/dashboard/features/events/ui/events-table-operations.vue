<script setup lang="ts">
import { flatOperations } from '~~/layers/dashboard/features/events/constants/helpers.js'

import type { EventOperation } from '~~/layers/dashboard/api/events.js'
import { getOperationColor } from '~~/layers/dashboard/features/events/composables/use-operation-color.js'
import type { EventOperationType } from '~/gql/graphql.js'

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
