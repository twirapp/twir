<script setup lang="ts">
import { CheckIcon, ChevronsUpDownIcon } from 'lucide-vue-next'
import { useField } from 'vee-validate'
import { computed } from 'vue'


import type { EventOperation } from '#layers/dashboard/api/events.ts'
import { getOperationColor } from '../composables/use-operation-color.ts'





import { flatOperations } from '~/features/events/constants/helpers.ts'
import { EventOperations } from '~/features/events/constants/operations.ts'
import { EventOperationType } from '~/gql/graphql.ts'
import { cn } from '~/lib/utils.ts'

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
	childrens: Array<{ name: string; value: EventOperationType }>
}>(([key, value]) => {
	if (value.childrens) {
		return {
			isGroup: true,
			name: value.name,
			childrens: Object.entries(value.childrens!).map(([childKey, childValue]) => ({
				name: childValue.name,
				value: Object.values(EventOperationType).find((v) => v === childKey)!,
			})),
		}
	}

	return {
		isGroup: false,
		name: value.name,
		value: Object.values(EventOperationType).find((v) => v === key)!,
		childrens: [],
	}
})
</script>

<template>
	<UiFormField :name="`operations.${currentOperationIndex}.type`">
		<UiFormItem class="flex flex-col">
			<UiFormLabel>{{ t('events.operations.name') }}</UiFormLabel>
			<UiFormControl>
				<UiPopover>
					<UiPopoverTrigger as-child>
						<UiFormControl>
							<UiButton
								type="button"
								variant="outline"
								role="combobox"
								:class="
									cn(
										'w-[300px] max-w-[300px] justify-between truncate',
										!currentOperation?.type && 'text-muted-foreground'
									)
								"
							>
								<div class="flex items-center gap-2">
									<div
										class="rounded-full size-3"
										:class="[getOperationColor(currentOperation?.type)]"
									></div>
									{{
										flatOperations[currentOperation?.type]
											? flatOperations[currentOperation?.type].name
											: 'Select...'
									}}
								</div>
								<ChevronsUpDownIcon class="ml-2 h-4 w-4 shrink-0 opacity-50" />
							</UiButton>
						</UiFormControl>
					</UiPopoverTrigger>
					<UiPopoverContent class="w-[400px] p-0">
						<UiCommand>
							<UiCommandInput placeholder="Search language..." />
							<UiCommandEmpty>Nothing found.</UiCommandEmpty>
							<UiCommandList>
								<template v-for="selectOption of typeSelectOptions">
									<UiCommandGroup
										v-if="selectOption.isGroup"
										:key="selectOption.name"
										:heading="selectOption.name"
									>
										<UiCommandItem
											v-for="operation of selectOption.childrens"
											:key="operation.name"
											:value="operation.value"
											@select="
												() => {
													if (!currentOperation) return
													currentOperation.type = operation.value
												}
											"
											class="flex items-center gap-2"
										>
											<div
												class="rounded-full size-3"
												:class="[getOperationColor(operation.value)]"
											></div>
											{{ operation.name }}
											<CheckIcon
												:class="
													cn(
														'ml-auto h-4 w-4',
														currentOperation?.type === operation.value ? 'opacity-100' : 'opacity-0'
													)
												"
											/>
										</UiCommandItem>
									</UiCommandGroup>

									<UiCommandGroup v-else>
										<UiCommandItem
											:key="selectOption.name!"
											:value="selectOption.value!"
											@select="
												() => {
													if (!currentOperation) return
													currentOperation.type = selectOption.value!
												}
											"
											class="flex items-center gap-2"
										>
											<div
												class="rounded-full size-3"
												:class="[getOperationColor(selectOption.value)]"
											></div>
											{{ selectOption.name }}
											<CheckIcon
												:class="
													cn(
														'ml-auto h-4 w-4',
														currentOperation?.type === selectOption.value
															? 'opacity-100'
															: 'opacity-0'
													)
												"
											/>
										</UiCommandItem>
									</UiCommandGroup>
								</template>
								<UiCommandGroup> </UiCommandGroup>
							</UiCommandList>
						</UiCommand>
					</UiPopoverContent>
				</UiPopover>
			</UiFormControl>
			<UiFormMessage />
		</UiFormItem>
	</UiFormField>
</template>
