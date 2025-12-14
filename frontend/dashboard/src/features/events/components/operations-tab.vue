<script setup lang="ts">
import { GripVerticalIcon, PlusIcon, TrashIcon } from 'lucide-vue-next'
import { useFieldArray } from 'vee-validate'
import { ref, useTemplateRef } from 'vue'
import { useDraggable } from 'vue-draggable-plus'
import { useI18n } from 'vue-i18n'

import OperationDetails from './operation.vue'

import type { EventOperation } from '@/api/events'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { flatOperations } from '@/features/events/constants/helpers'
import { EventOperationType } from '@/gql/graphql'
import { getOperationColor } from '@/features/events/composables/use-operation-color.ts'

const { t } = useI18n()

const selectedOperation = ref(0)

const operations = useFieldArray<Omit<EventOperation, 'id'>>('operations')

function addOperation() {
	if (operations.fields.value.length >= 10) {
		return
	}

	operations.insert(operations.fields?.value?.length ?? 0, {
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

const draggableRef = useTemplateRef('draggableRef')
useDraggable(draggableRef, operations.fields, {
	animation: 150,
	handle: '.drag-handle',
	onUpdate(event) {
		if (event.oldIndex !== undefined && event.newIndex !== undefined) {
			operations.move(event.oldIndex, event.newIndex)
			selectedOperation.value = event.newIndex
		}
	},
})
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle>{{ t('events.operations.name') }}</CardTitle>
		</CardHeader>
		<CardContent>
			<div class="flex flex-col lg:flex-row gap-4">
				<div
					class="basis-1/3 h-full w-full lg:w-[30%] lg:max-w-[30%] lg:min-w-[30%] flex flex-col gap-2 flex-wrap items-center rounded-lg border p-3 shadow-xs"
				>
					<Button
						type="button"
						:disabled="operations.fields.value.length >= 10"
						class="flex items-center gap-2 w-full"
						variant="outline"
						@click="addOperation"
					>
						<template v-if="operations.fields.value.length < 10">
							<PlusIcon class="size-4" />
							Create new
						</template>
						<template v-else> Maximum limit reached </template>
					</Button>
					<div ref="draggableRef" class="w-full flex flex-col gap-2">
						<template v-if="operations.fields.value.length">
							<div
								v-for="(operation, operationIndex) in operations.fields.value"
								:key="operationIndex"
								class="w-full rounded-lg border py-1 px-2 cursor-pointer items-center flex flex-row justify-between select-none"
								@click="selectOperation(operationIndex)"
							>
								<div class="flex gap-2 items-center">
									<GripVerticalIcon class="min-size-4 cursor-grab drag-handle" />
									<div
										class="rounded-full size-3"
										:class="[getOperationColor(operation.value?.type)]"
									></div>
									<span v-if="operation.value" class="truncate">
										{{
											flatOperations[operation.value?.type]?.name?.slice(0, 30) ??
											'Unknown Operation'
										}}
									</span>
								</div>
								<div class="flex gap-2 items-center">
									<Button
										type="button"
										class="flex items-center"
										size="sm"
										variant="ghost"
										@click.stop="() => removeOperation(operationIndex)"
									>
										<TrashIcon class="size-4" />
									</Button>
								</div>
							</div>
						</template>
					</div>
				</div>

				<div class="w-full h-full">
					<OperationDetails :operation-index="selectedOperation" />
				</div>
			</div>
		</CardContent>
	</Card>
</template>
