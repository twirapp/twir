<script setup lang="ts">
import { CheckIcon, ChevronsUpDownIcon, PlusIcon } from 'lucide-vue-next'
import { useField, useFieldArray } from 'vee-validate'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import OperationFilter from './operation-filter.vue'

import type { EventFilter, EventOperation } from '@/api/events'

import { flatOperations } from '@/components/events/helpers'
import { EventOperations } from '@/components/events/operations'
import { Button } from '@/components/ui/button'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
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
	Popover,
	PopoverContent,
	PopoverTrigger,
} from '@/components/ui/popover'
import { Separator } from '@/components/ui/separator'
import { Switch } from '@/components/ui/switch'
import VariableInput from '@/components/variable-input.vue'
import { EventOperationType } from '@/gql/graphql'
import { cn } from '@/lib/utils'

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

const typeSelectOptions = Object.entries(EventOperations).map<{
	isGroup: boolean
	name: string
	value?: EventOperationType
	childrens: Array<{ name: string, value: EventOperationType }>
}>(([key, value]) => {
	if (value.childrens) {
		return {
			isGroup: true,
			name: value.name,
			childrens: Object.entries(value.childrens!).map(([childKey, childValue]) => ({
				name: childValue.name,
				value: Object.values(EventOperationType).find(v => v === childKey)!,
			})),
		}
	}

	return {
		isGroup: false,
		name: value.name,
		value: Object.values(EventOperationType).find(v => v === key)!,
		childrens: [],
	}
})

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
			<FormField
				:name="`operations.${operationIndex}.type`"
			>
				<FormItem class="flex flex-col">
					<FormLabel>{{ t('events.operations.name') }}</FormLabel>
					<FormControl>
						<Popover>
							<PopoverTrigger as-child>
								<FormControl>
									<Button
										variant="outline"
										role="combobox"
										:class="cn('w-[300px] max-w-[300px] justify-between truncate', !currentOperation?.type && 'text-muted-foreground')"
									>
										{{ flatOperations[currentOperation?.type] ? flatOperations[currentOperation?.type].name : 'Select...' }}
										<ChevronsUpDownIcon class="ml-2 h-4 w-4 shrink-0 opacity-50" />
									</Button>
								</FormControl>
							</PopoverTrigger>
							<PopoverContent class="w-[400px] p-0">
								<Command>
									<CommandInput placeholder="Search language..." />
									<CommandEmpty>Nothing found.</CommandEmpty>
									<CommandList>
										<template v-for="selectOption of typeSelectOptions">
											<CommandGroup
												v-if="selectOption.isGroup"
												:key="selectOption.name"
												:heading="selectOption.name"
											>
												<CommandItem
													v-for="operation of selectOption.childrens"
													:key="operation.name"
													:value="operation.value"
													@select="() => {
														if (!currentOperation) return;
														currentOperation.type = operation.value;
													}"
												>
													{{ operation.name }}
													<CheckIcon
														:class="cn('ml-auto h-4 w-4', currentOperation?.type === operation.value ? 'opacity-100' : 'opacity-0')"
													/>
												</CommandItem>
											</CommandGroup>

											<CommandItem
												v-else
												:key="selectOption.name!"
												:value="selectOption.value!"
												@select="() => {
													if (!currentOperation) return;
													currentOperation.type = selectOption.value!;
												}"
											>
												{{ selectOption.name }}
												<CheckIcon :class="cn('ml-auto h-4 w-4', currentOperation?.type === selectOption.value ? 'opacity-100' : 'opacity-0')" />
											</CommandItem>
										</template>
										<CommandGroup>
										</CommandGroup>
									</CommandList>
								</Command>
							</PopoverContent>
						</Popover>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

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
								:placeholder="t('events.delayPlaceholder')"
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
					<Button variant="outline" size="sm" @click="() => onAddFilter()">
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
			Create operation first.
		</div>
	</div>
</template>
