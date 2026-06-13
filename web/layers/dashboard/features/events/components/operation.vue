<script lang="ts" setup>
import { PlusIcon } from 'lucide-vue-next'
import { useField, useFieldArray } from 'vee-validate'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import OperationFilter from './filter.vue'
import { useCommandsApi } from '@/api/commands/commands.ts'
import { useVariablesApi } from '@/api/variables.ts'
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
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import VariableInput from '@/components/variable-input.vue'
import OperationActionSelector from '@/features/events/components/operation-action-selector.vue'
import OperationInputAlert from '@/features/events/components/operation-input-alert.vue'
import OperationInputObsSelector from '@/features/events/components/operation-input-obs-selector.vue'
import { flatOperations } from '@/features/events/constants/helpers'
import { EventOperationType } from '@/gql/graphql'

import type { EventFilter, EventOperation } from '@/api/events'

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
		<div
			v-if="currentOperation"
			class="space-y-4"
		>
			<OperationActionSelector :current-operation-index="operationIndex" />

			<div v-if="flatOperations[currentOperation?.type]?.haveInput">
				<FormField
					v-slot="{ componentField }"
					:name="`operations[${operationIndex}].input`"
				>
					<FormItem>
						<FormLabel>
							{{
								t(
									flatOperations[currentOperation?.type].inputKeyTranslatePath ??
										'events.operations.inputs.default'
								)
							}}
						</FormLabel>
						<FormControl>
							<Input
								v-if="currentOperation?.type === EventOperationType.SendHttpRequest"
								placeholder="https://example.com/webhook"
								type="url"
								v-bind="componentField"
							/>
							<VariableInput
								v-else
								input-type="textarea"
								v-bind="componentField"
							/>
						</FormControl>
						<FormDescription v-if="currentOperation?.type === EventOperationType.SendHttpRequest">
							A POST request will be sent to this URL with the event data as JSON body. The event
							type will be in the <code>X-Twir-Event-Type</code> header.
						</FormDescription>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
				<FormField
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.delay`"
				>
					<FormItem>
						<FormLabel>{{ t('events.delay') }}</FormLabel>
						<FormControl>
							<Input
								min="0"
								type="number"
								v-bind="componentField"
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
								:placeholder="t('events.repeatPlaceholder')"
								min="0"
								type="number"
								v-bind="componentField"
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
					<FormItem
						class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-xs"
					>
						<div class="space-y-0.5">
							<FormLabel>Use announce</FormLabel>
						</div>
						<FormControl>
							<Switch
								:model-value="value"
								@update:model-value="handleChange"
							/>
						</FormControl>
					</FormItem>
				</FormField>
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
				<FormField
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.timeoutTime`"
				>
					<FormItem>
						<FormLabel>{{ t('events.operations.banTime') }}</FormLabel>
						<FormControl>
							<Input
								:placeholder="t('events.operations.banTime')"
								min="0"
								type="number"
								v-bind="componentField"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<div
				v-if="flatOperations[currentOperation?.type]?.additionalValues?.includes('timeoutMessage')"
			>
				<FormField
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.timeoutMessage`"
				>
					<FormItem>
						<FormLabel>{{ t('events.operations.banMessage') }}</FormLabel>
						<FormControl>
							<VariableInput
								class="relative"
								input-type="textarea"
								v-bind="componentField"
							/>
						</FormControl>
						<FormMessage />
					</FormItem>
				</FormField>
			</div>

			<div v-if="flatOperations[currentOperation?.type]?.additionalValues?.includes('target')">
				<FormField
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
					<FormItem>
						<FormLabel>Variable</FormLabel>

						<Select v-bind="componentField">
							<FormControl>
								<SelectTrigger>
									<SelectValue placeholder="Select a variable" />
								</SelectTrigger>
							</FormControl>
							<SelectContent>
								<SelectGroup>
									<SelectItem
										v-for="variable of customVariables"
										:value="variable.id"
									>
										{{ variable.name }}
									</SelectItem>
								</SelectGroup>
							</SelectContent>
						</Select>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField
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
					<FormItem>
						<FormLabel>Variable</FormLabel>

						<Select v-bind="componentField">
							<FormControl>
								<SelectTrigger>
									<SelectValue placeholder="Select a command" />
								</SelectTrigger>
							</FormControl>
							<SelectContent>
								<SelectGroup>
									<SelectItem
										v-for="command of commands?.commands"
										:value="command.id"
									>
										{{ command.name }}
									</SelectItem>
								</SelectGroup>
							</SelectContent>
						</Select>
						<FormMessage />
					</FormItem>
				</FormField>

				<FormField
					v-else
					v-slot="{ componentField }"
					:name="`operations.${operationIndex}.target`"
				>
					<FormItem>
						<FormLabel>{{ t('events.target') }}</FormLabel>
						<FormControl>
							<Input
								:placeholder="t('events.targetPlaceholder')"
								v-bind="componentField"
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
				<FormItem
					class="flex flex-row items-center justify-between rounded-lg border p-3 shadow-xs"
				>
					<div class="space-y-0.5">
						<FormLabel>{{ t('sharedTexts.enabled') }}?</FormLabel>
					</div>
					<FormControl>
						<Switch
							:model-value="value"
							@update:model-value="handleChange"
						/>
					</FormControl>
				</FormItem>
			</FormField>

			<Separator />

			<div class="mt-6">
				<div class="mb-4 flex items-center justify-between">
					<h3 class="text-lg font-medium">
						{{ t('events.operations.filters.label') }}
					</h3>
					<Button
						size="sm"
						type="button"
						variant="outline"
						@click="() => onAddFilter()"
					>
						<PlusIcon class="mr-2 h-4 w-4" />
						{{ t('sharedTexts.create') }}
					</Button>
				</div>

				<div
					v-if="currentOperation.filters?.length === 0"
					class="rounded-md border p-4 text-center"
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
			class="flex grow flex-col items-center justify-center rounded-md border-2 border-dashed p-4"
		>
			Create operation.
		</div>
	</div>
</template>
