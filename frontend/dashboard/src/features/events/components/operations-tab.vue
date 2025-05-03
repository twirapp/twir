<script setup lang="ts">
import { GripVerticalIcon, PlusIcon, TrashIcon } from 'lucide-vue-next'
import { useFieldArray } from 'vee-validate'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'

import OperationDetails from './operation-details.vue'

import type { EventOperation } from '@/api/events'

import { flatOperations } from '@/components/events/helpers'
import { Button } from '@/components/ui/button'
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card'
import { EventOperationType } from '@/gql/graphql'

const { t } = useI18n()

const selectedOperation = ref(0)

const operations = useFieldArray<Omit<EventOperation, 'id'>>('operations')

function addOperation() {
	if (operations.fields.value.length >= 10) {
		return
	}

	operations.insert(operations.fields.value.length - 1, {
		type: EventOperationType.SendMessage,
		delay: 0,
		enabled: true,
		filters: [],
		repeat: 0,
		timeoutTime: 0,
		input: '',
		target: '',
		timeoutMessage: '',
		useAnnounce: false,
	})

	selectedOperation.value = operations.fields.value.length - 1
}

function selectOperation(operationIndex: number) {
	selectedOperation.value = operationIndex
}

function removeOperation(operationIndex: number) {
	operations.remove(operationIndex)
	selectedOperation.value = 0
}
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle>{{ t('events.operations') }}</CardTitle>
			<CardDescription>{{ t('events.operationsDescription') }}</CardDescription>
		</CardHeader>
		<CardContent>
			<div class="flex flex-col lg:flex-row gap-4">
				<div
					class="basis-1/3 h-full w-[30%] max-w-[30%] min-w-[30%] flex flex-col gap-2 flex-wrap items-center rounded-lg border p-3 shadow-sm"
				>
					<template v-if="operations.fields.value.length">
						<div
							v-for="(operation, operationIndex) in operations.fields.value" :key="operationIndex"
							class="w-full rounded-lg border py-1 px-2 cursor-pointer items-center flex flex-row justify-between"
							:class="{
								'outline outline-1 outline-zinc-700': operationIndex === selectedOperation,
							}"
							@click="selectOperation(operationIndex)"
						>
							<GripVerticalIcon class="min-size-4 cursor-grab" />
							<span v-if="operation.value" class="truncate">
								{{ flatOperations[operation.value.type]?.name ?? 'Unknown Operation' }}
							</span>
							<Button
								class="flex items-center"
								size="sm"
								variant="ghost"
								@click.stop="() => removeOperation(operationIndex)"
							>
								<TrashIcon class="size-4" />
							</Button>
						</div>
					</template>

					<Button
						:disabled="operations.fields.value.length >= 10"
						class="flex items-center gap-2 w-full"
						variant="outline"
						@click="addOperation"
					>
						<template v-if="operations.fields.value.length < 10">
							<PlusIcon class="size-4" />
							Create new
						</template>
						<template v-else>
							Maximum limit reached
						</template>
					</Button>
				</div>

				<div class="w-full h-full">
					<OperationDetails :operation-index="selectedOperation" />
				</div>
			</div>
		</CardContent>
	</Card>
</template>
