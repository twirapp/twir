<script setup lang="ts">
import { CheckIcon, ChevronsUpDownIcon } from 'lucide-vue-next'
import { useField } from 'vee-validate'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import type { EventOperation } from '@/api/events.ts'

import { Button } from '@/components/ui/button'
import {
	Command,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from '@/components/ui/command'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { flatOperations } from '@/features/events/constants/helpers.ts'
import { EventOperations } from '@/features/events/constants/operations.ts'
import { EventOperationType } from '@/gql/graphql.ts'
import { cn } from '@/lib/utils.ts'

const props = defineProps<{
	currentOperationIndex: number
}>()
const { t } = useI18n()

const currentOperationPath = computed(() => `operations.${props.currentOperationIndex}`)
const { value: currentOperation } = useField<Omit<EventOperation, 'id'>>(currentOperationPath)

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
</script>

<template>
	<FormField
		:name="`operations.${currentOperationIndex}.type`"
	>
		<FormItem class="flex flex-col">
			<FormLabel>{{ t('events.operations.name') }}</FormLabel>
			<FormControl>
				<Popover>
					<PopoverTrigger as-child>
						<FormControl>
							<Button
								type="button"
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
</template>
