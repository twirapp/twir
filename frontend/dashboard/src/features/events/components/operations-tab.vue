<script setup lang="ts">
import { ArrowDown, ArrowUp, PlusIcon, Trash2 } from 'lucide-vue-next'
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'

import DraggableOperations from './draggable-operations.vue'
import OperationDetails from './operation-details.vue'

import { Button } from '@/components/ui/button'
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from '@/components/ui/card'
import {
	Tabs,
	TabsContent,
	TabsList,
	TabsTrigger,
} from '@/components/ui/tabs'

const { t } = useI18n()

// Operations management
const selectedOperationTab = ref(0)
const currentOperation = computed(() => {
	const operations = props.form.values.operations
	return operations.length > 0 ? operations[selectedOperationTab.value] : null
})

function addOperation() {
	const operations = [...props.form.values.operations]
	operations.push({
		type: '',
		input: '',
		delay: 0,
		repeat: 0,
		useAnnounce: false,
		timeoutTime: 0,
		timeoutMessage: '',
		target: '',
		enabled: true,
		filters: [],
	})
	props.form.setFieldValue('operations', operations)
	selectedOperationTab.value = operations.length - 1
}

function removeOperation(index: number) {
	const operations = [...props.form.values.operations]
	operations.splice(index, 1)
	props.form.setFieldValue('operations', operations)

	if (selectedOperationTab.value >= operations.length) {
		selectedOperationTab.value = Math.max(0, operations.length - 1)
	}
}

function addFilter(operationIndex: number) {
	const operations = [...props.form.values.operations]
	operations[operationIndex].filters.push({
		type: 'EQUALS',
		left: '',
		right: '',
	})
	props.form.setFieldValue('operations', operations)
}

function removeFilter(operationIndex: number, filterIndex: number) {
	const operations = [...props.form.values.operations]
	operations[operationIndex].filters.splice(filterIndex, 1)
	props.form.setFieldValue('operations', operations)
}

function moveOperationUp(index: number) {
	if (index <= 0) return // Уже первый элемент

	const operations = [...props.form.values.operations]
	const temp = operations[index]
	operations[index] = operations[index - 1]
	operations[index - 1] = temp

	props.form.setFieldValue('operations', operations)
	selectedOperationTab.value = index - 1
}

function moveOperationDown(index: number) {
	const operations = [...props.form.values.operations]
	if (index >= operations.length - 1) return // Уже последний элемент

	const temp = operations[index]
	operations[index] = operations[index + 1]
	operations[index + 1] = temp

	props.form.setFieldValue('operations', operations)
	selectedOperationTab.value = index + 1
}
</script>

<template>
	<Card>
		<CardHeader>
			<CardTitle>{{ t('events.operations') }}</CardTitle>
			<CardDescription>{{ t('events.operationsDescription') }}</CardDescription>
		</CardHeader>
		<CardContent>
			<div v-if="form.values.operations.length === 0" class="text-center py-8">
				<p class="text-muted-foreground mb-4">
					{{ t('events.noOperations') }}
				</p>
				<Button type="button" @click="addOperation">
					<PlusIcon class="mr-2 h-4 w-4" />
					{{ t('events.addOperation') }}
				</Button>
			</div>

			<div v-else>
				<DraggableOperations
					:form="form"
					:selected-tab="selectedOperationTab"
					:on-tab-change="(index) => selectedOperationTab = index"
				/>
				<Tabs v-model="selectedOperationTab">
					<div class="flex justify-between items-center mb-4">
						<TabsList>
							<TabsTrigger
								v-for="(operation, index) in form.values.operations"
								:key="index"
								:value="index"
							>
								{{ t('events.operation') }} {{ index + 1 }}
							</TabsTrigger>
						</TabsList>

						<div class="flex gap-2">
							<Button
								type="button"
								variant="outline"
								size="sm"
								:disabled="selectedOperationTab <= 0"
								@click="moveOperationUp(selectedOperationTab)"
							>
								<ArrowUp class="h-4 w-4" />
							</Button>
							<Button
								type="button"
								variant="outline"
								size="sm"
								:disabled="selectedOperationTab >= form.values.operations.length - 1"
								@click="moveOperationDown(selectedOperationTab)"
							>
								<ArrowDown class="h-4 w-4" />
							</Button>
							<Button type="button" variant="outline" size="sm" @click="addOperation">
								<PlusIcon class="h-4 w-4" />
							</Button>
							<Button
								v-if="form.values.operations.length > 0"
								type="button"
								variant="destructive"
								size="sm"
								@click="removeOperation(selectedOperationTab)"
							>
								<Trash2 class="h-4 w-4" />
							</Button>
						</div>
					</div>

					<div v-for="(operation, opIndex) in form.values.operations" :key="opIndex">
						<TabsContent :value="opIndex">
							<OperationDetails
								:operation-index="opIndex"
								:operation="operation"
								:on-add-filter="addFilter"
								:on-remove-filter="removeFilter"
							/>
						</TabsContent>
					</div>
				</Tabs>
			</div>
		</CardContent>
	</Card>
</template>
