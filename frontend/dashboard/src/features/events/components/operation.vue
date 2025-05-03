<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { useField, useFieldArray } from 'vee-validate'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import OperationFilter from './filter.vue'

import type { EventFilter, EventOperation } from '@/api/events'

import { Button } from '@/components/ui/button'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import VariableInput from '@/components/variable-input.vue'
import OperationActionSelector from '@/features/events/components/operation-action-selector.vue'
import OperationInputAlert from '@/features/events/components/operation-input-alert.vue'
import OperationInputObsSelector
	from '@/features/events/components/operation-input-obs-selector.vue'
import { flatOperations } from '@/features/events/constants/helpers'
import { EventOperationType } from '@/gql/graphql'

const props = withDefaults(defineProps<{
	operationIndex: number
}>(), {
	operationIndex: 0,
})

const currentOperationPath = computed(() => `operations.${props.operationIndex}`)
const { value: currentOperation } = useField<Omit<EventOperation, 'id'> | undefined>(currentOperationPath)

const currentOperationFiltersPath = computed(() => `${currentOperationPath.value}.filters`)
const { fields: filters, insert: insertFilter, remove: removeFilter, replace: updateFilters } = useFieldArray<EventFilter | undefined>(currentOperationFiltersPath)

const { t } = useI18n()

function onAddFilter() {
	insertFilter(filters.value.length, {
		id: '',
		type: 'EQUALS',
		left: '',
		right: '',
	})
}

function onRemoveFilter(filterIndex: number) {
	if (filterIndex === 0) {
		updateFilters([])
	} else {
		removeFilter(filterIndex)
	}
}
</script>

<template>
	<div class="min-h-[50dvh]">
		<div v-if="currentOperation" class="space-y-4">
			<OperationActionSelector
				:current-operation-index="operationIndex"
			/>

			<div v-if="flatOperations[currentOperation?.type]?.haveInput">
				<FormField
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.input`"
				>
					<FormItem>
						<FormLabel>
							{{
								t(flatOperations[currentOperation?.type].inputKeyTranslatePath ?? 'events.operations.inputs.default')
							}}
						</FormLabel>
						<FormControl>
							<VariableInput v-bind="componentField" input-type="textarea" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<FormField
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.delay`"
				>
					<FormItem>
						<FormLabel>{{ t('events.delay') }}</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								type="number"
								min="0"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.repeat`"
				>
					<FormItem>
						<FormLabel>{{ t('events.repeat') }}</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								type="number"
								min="0"
								:placeholder="t('events.repeatPlaceholder')"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<div v-if="flatOperations[currentOperation?.type]?.additionalValues?.includes('useAnnounce')">
				<FormField
					v-slot="{ value, handleChange }"
					:name="`operations.${operationIndex}.useAnnounce`"
				>
					<FormItem class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
						<div class="space-y-0.5">
							<FormLabel>Use announce</FormLabel>
						</div>
						<FormControl>
							<Switch
								:checked="value"
								@update:checked="handleChange"
							/>
						</FormControl>
					</FormItem>
				</FormField>
			</div>

			<div v-if="currentOperation?.type === EventOperationType.TriggerAlert">
				<OperationInputAlert :operation-index="operationIndex" />
			</div>

			<div
				v-if="[
					EventOperationType.ObsChangeScene,
					EventOperationType.ObsDecreaseAudioVolume,
					EventOperationType.ObsDisableAudio,
					EventOperationType.ObsEnableAudio,
					EventOperationType.ObsIncreaseAudioVolume,
					EventOperationType.ObsSetAudioVolume,
					EventOperationType.ObsStartStream,
					EventOperationType.ObsStopStream,
					EventOperationType.ObsToggleAudio,
					EventOperationType.ObsToggleSource,
				].includes(currentOperation.type)"
			>
				<OperationInputObsSelector
					:operation-index="operationIndex"
				/>
			</div>

			<div v-if="flatOperations[currentOperation?.type]?.additionalValues?.includes('timeoutTime')">
				<FormField
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.timeoutTime`"
				>
					<FormItem>
						<FormLabel>{{ t('events.operations.banTime') }}</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								type="number"
								min="0"
								:placeholder="t('events.operations.banTime')"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<div v-if="flatOperations[currentOperation?.type]?.additionalValues?.includes('timeoutMessage')">
				<FormField
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.timeoutMessage`"
				>
					<FormItem>
						<FormLabel>{{ t('events.operations.banMessage') }}</FormLabel>
						<FormControl>
							<VariableInput v-bind="componentField" input-type="textarea" />
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<div v-if="flatOperations[currentOperation?.type]?.additionalValues?.includes('target')">
				<FormField
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.target`"
				>
					<FormItem>
						<FormLabel>{{ t('events.target') }}</FormLabel>
						<FormControl>
							<Input
								v-bind="componentField"
								:placeholder="t('events.targetPlaceholder')"
							/>
						</FormControl>
						<FormDescription>{{ t('events.targetDescription') }}</FormDescription>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<FormField
				v-slot="{ value, handleChange }"
				:name="`operations.${operationIndex}.enabled`"
			>
				<FormItem class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm">
					<div class="space-y-0.5">
						<FormLabel>{{ t('sharedTexts.enabled') }}?</FormLabel>
					</div>
					<FormControl>
						<Switch
							:checked="value"
							@update:checked="handleChange"
						/>
					</FormControl>
				</FormItem>
			</FormField>

			<Separator />

			<div class="mt-6">
				<div class="flex justify-between items-center mb-4">
					<h3 class="text-lg font-medium">
						{{ t('events.operations.filters.label') }}
					</h3>
					<Button type="button" variant="outline" size="sm" @click="() => onAddFilter()">
						<PlusIcon class="h-4 w-4 mr-2" />
						{{ t('sharedTexts.create') }}
					</Button>
				</div>

				<div v-if="currentOperation.filters?.length === 0" class="text-center p-4 border rounded-md">
					<p class="text-muted-foreground">
						{{ t('events.operations.filters.description') }}
					</p>
				</div>

				<template v-else>
					<OperationFilter
						v-for="(_, filterIndex) in currentOperation.filters"
						:key="filterIndex"
						:operation-index="operationIndex"
						:filter-index="filterIndex"
						:on-remove="onRemoveFilter"
					/>
				</template>
			</div>
		</div>

		<div
			v-else
			class="flex flex-col flex-grow items-center justify-center border-2 rounded-md p-4 border-dashed"
		>
			Create operation.
		</div>
	</div>
</template>
