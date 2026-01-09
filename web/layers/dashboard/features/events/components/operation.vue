<script setup lang="ts">
import { PlusIcon } from 'lucide-vue-next'
import { useField, useFieldArray } from 'vee-validate'
import { computed } from 'vue'


import OperationFilter from './filter.vue'

import type { EventFilter, EventOperation } from '#layers/dashboard/api/events'







import VariableInput from '#layers/dashboard/components/variable-input.vue'
import OperationActionSelector from '~/features/events/components/operation-action-selector.vue'
import OperationInputAlert from '~/features/events/components/operation-input-alert.vue'
import OperationInputObsSelector from '~/features/events/components/operation-input-obs-selector.vue'
import { flatOperations } from '~/features/events/constants/helpers'
import { EventOperationType } from '~/gql/graphql'
import { useVariablesApi } from '#layers/dashboard/api/variables.ts'
import { useCommandsApi } from '#layers/dashboard/api/commands/commands.ts'

const props = withDefaults(
	defineProps<{
		operationIndex: number
	}>(),
	{
		operationIndex: 0,
	}
)

const { customVariables } = useVariablesApi()
const commandsApi = useCommandsApi()
const { data: commands } = commandsApi.useQueryCommands()

const currentOperationPath = computed(() => `operations.${props.operationIndex}`)
const { value: currentOperation } = useField<Omit<EventOperation, 'id'> | undefined>(
	currentOperationPath
)

const currentOperationFiltersPath = computed(() => `${currentOperationPath.value}.filters`)
const {
	fields: filters,
	insert: insertFilter,
	remove: removeFilter,
	replace: updateFilters,
} = useFieldArray<EventFilter | undefined>(currentOperationFiltersPath)

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
			<OperationActionSelector :current-operation-index="operationIndex" />

			<div v-if="flatOperations[currentOperation?.type]?.haveInput">
				<UiFormField v-slot="{ componentField }" :name="`operations[${operationIndex}].input`">
					<UiFormItem>
						<UiFormLabel>
							{{
								t(
									flatOperations[currentOperation?.type].inputKeyTranslatePath ??
										'events.operations.inputs.default'
								)
							}}
						</UiFormLabel>
						<UiFormControl>
							<VariableInput v-bind="componentField" input-type="textarea" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>

			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<UiFormField v-slot="{ componentField }" :name="`operations.${operationIndex}.delay`">
					<UiFormItem>
						<UiFormLabel>{{ t('events.delay') }}</UiFormLabel>
						<UiFormControl>
							<UiInput v-bind="componentField" type="number" min="0" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ componentField }" :name="`operations.${operationIndex}.repeat`">
					<UiFormItem>
						<UiFormLabel>{{ t('events.repeat') }}</UiFormLabel>
						<UiFormControl>
							<UiInput
								v-bind="componentField"
								type="number"
								min="0"
								:placeholder="t('events.repeatPlaceholder')"
							/>
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>

			<div v-if="flatOperations[currentOperation?.type]?.additionalValues?.includes('useAnnounce')">
				<UiFormField
					v-slot="{ value, handleChange }"
					:name="`operations.${operationIndex}.useAnnounce`"
				>
					<UiFormItem
						class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-xs"
					>
						<div class="space-y-0.5">
							<UiFormLabel>Use announce</UiFormLabel>
						</div>
						<UiFormControl>
							<UiSwitch :model-value="value" @update:model-value="handleChange" />
						</UiFormControl>
					</UiFormItem>
				</UiFormField>
			</div>

			<div v-if="currentOperation?.type === EventOperationType.TriggerAlert">
				<OperationInputAlert :operation-index="operationIndex" />
			</div>

			<div
				v-if="
					[
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
					].includes(currentOperation.type)
				"
			>
				<OperationInputObsSelector :operation-index="operationIndex" />
			</div>

			<div v-if="flatOperations[currentOperation?.type]?.additionalValues?.includes('timeoutTime')">
				<UiFormField v-slot="{ componentField }" :name="`operations.${operationIndex}.timeoutTime`">
					<UiFormItem>
						<UiFormLabel>{{ t('events.operations.banTime') }}</UiFormLabel>
						<UiFormControl>
							<UiInput
								v-bind="componentField"
								type="number"
								min="0"
								:placeholder="t('events.operations.banTime')"
							/>
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>

			<div
				v-if="flatOperations[currentOperation?.type]?.additionalValues?.includes('timeoutMessage')"
			>
				<UiFormField
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.timeoutMessage`"
				>
					<UiFormItem>
						<UiFormLabel>{{ t('events.operations.banMessage') }}</UiFormLabel>
						<UiFormControl>
							<VariableInput v-bind="componentField" input-type="textarea" class="relative" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>

			<div v-if="flatOperations[currentOperation?.type]?.additionalValues?.includes('target')">
				<UiFormField
					v-if="
						[
							EventOperationType.ChangeVariable,
							EventOperationType.DecrementVariable,
							EventOperationType.IncrementVariable,
						].includes(currentOperation.type)
					"
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.target`"
				>
					<UiFormItem>
						<UiFormLabel>Variable</UiFormLabel>

						<UiSelect v-bind="componentField">
							<UiFormControl>
								<UiSelectTrigger>
									<UiSelectValue placeholder="Select a variable" />
								</UiSelectTrigger>
							</UiFormControl>
							<UiSelectContent>
								<UiSelectGroup>
									<UiSelectItem v-for="variable of customVariables" :value="variable.id">
										{{ variable.name }}
									</UiSelectItem>
								</UiSelectGroup>
							</UiSelectContent>
						</UiSelect>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField
					v-else-if="
						[
							EventOperationType.AllowCommandToUser,
							EventOperationType.RemoveAllowCommandToUser,
							EventOperationType.DenyCommandToUser,
							EventOperationType.RemoveDenyCommandToUser,
						].includes(currentOperation.type)
					"
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.target`"
				>
					<UiFormItem>
						<UiFormLabel>Variable</UiFormLabel>

						<UiSelect v-bind="componentField">
							<UiFormControl>
								<UiSelectTrigger>
									<UiSelectValue placeholder="Select a command" />
								</UiSelectTrigger>
							</UiFormControl>
							<UiSelectContent>
								<UiSelectGroup>
									<UiSelectItem v-for="command of commands?.commands" :value="command.id">
										{{ command.name }}
									</UiSelectItem>
								</UiSelectGroup>
							</UiSelectContent>
						</UiSelect>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField v-else v-slot="{ componentField }" :name="`operations.${operationIndex}.target`">
					<UiFormItem>
						<UiFormLabel>{{ t('events.target') }}</UiFormLabel>
						<UiFormControl>
							<UiInput v-bind="componentField" :placeholder="t('events.targetPlaceholder')" />
						</UiFormControl>
						<UiFormDescription>{{ t('events.targetDescription') }}</UiFormDescription>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
			</div>

			<UiFormField v-slot="{ value, handleChange }" :name="`operations.${operationIndex}.enabled`">
				<UiFormItem
					class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-xs"
				>
					<div class="space-y-0.5">
						<UiFormLabel>{{ t('sharedTexts.enabled') }}?</UiFormLabel>
					</div>
					<UiFormControl>
						<UiSwitch :model-value="value" @update:model-value="handleChange" />
					</UiFormControl>
				</UiFormItem>
			</UiFormField>

			<UiSeparator />

			<div class="mt-6">
				<div class="flex justify-between items-center mb-4">
					<h3 class="text-lg font-medium">
						{{ t('events.operations.filters.label') }}
					</h3>
					<UiButton type="button" variant="outline" size="sm" @click="() => onAddFilter()">
						<PlusIcon class="h-4 w-4 mr-2" />
						{{ t('sharedTexts.create') }}
					</UiButton>
				</div>

				<div
					v-if="currentOperation.filters?.length === 0"
					class="text-center p-4 border rounded-md"
				>
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
			class="flex flex-col grow items-center justify-center border-2 rounded-md p-4 border-dashed"
		>
			Create operation.
		</div>
	</div>
</template>
