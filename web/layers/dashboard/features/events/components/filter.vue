<script setup lang="ts">
import { Trash2 } from 'lucide-vue-next'





import VariableInput from '#layers/dashboard/components/variable-input.vue'

defineProps<{
	operationIndex: number
	filterIndex: number
	onRemove: (operationIndex: number, filterIndex: number) => void
}>()

const { t } = useI18n()
</script>

<template>
	<div class="border rounded-md p-4 mb-4">
		<div class="flex justify-between items-center mb-4">
			<h4 class="font-medium">
				{{ t('events.operations.filters.label') }} {{ filterIndex + 1 }}
			</h4>
			<UiButton
				type="button"
				variant="destructive"
				size="sm"
				@click="onRemove(operationIndex, filterIndex)"
			>
				<Trash2 class="h-4 w-4" />
			</UiButton>
		</div>

		<div class="space-y-4">
			<UiFormField
				v-slot="{ componentField }"
				:name="`operations.${operationIndex}.filters.${filterIndex}.left`"
			>
				<UiFormItem>
					<UiFormLabel>{{ t('events.operations.filters.placeholderLeft') }}</UiFormLabel>
					<UiFormControl>
						<VariableInput v-bind="componentField" />
					</UiFormControl>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiFormField
				v-slot="{ componentField }"
				:name="`operations.${operationIndex}.filters.${filterIndex}.type`"
			>
				<UiFormItem>
					<UiFormControl>
						<UiSelect
							v-bind="componentField"
							:placeholder="t('events.selectFilterType')"
						>
							<UiSelectTrigger>
								<UiSelectValue :placeholder="t('events.selectFilterType')" />
							</UiSelectTrigger>
							<UiSelectContent>
								<UiSelectItem value="EQUALS">
									=
								</UiSelectItem>
								<UiSelectItem value="NOT_EQUALS">
									!=
								</UiSelectItem>
								<UiSelectItem value="CONTAINS">
									contains
								</UiSelectItem>
								<UiSelectItem value="NOT_CONTAINS">
									not contains
								</UiSelectItem>
								<UiSelectItem value="STARTS_WITH">
									starts with
								</UiSelectItem>
								<UiSelectItem value="ENDS_WITH">
									ends with
								</UiSelectItem>
								<UiSelectItem value="GREATER_THAN">
									>
								</UiSelectItem>
								<UiSelectItem value="LESS_THAN">
									{{ '<' }}
								</UiSelectItem>
								<UiSelectItem value="GREATER_THAN_OR_EQUALS">
									>=
								</UiSelectItem>
								<UiSelectItem value="LESS_THAN_OR_EQUALS">
									{{ '<=' }}
								</UiSelectItem>
							</UiSelectContent>
						</UiSelect>
					</UiFormControl>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiFormField
				v-slot="{ componentField }"
				:name="`operations.${operationIndex}.filters.${filterIndex}.right`"
			>
				<UiFormItem>
					<UiFormLabel>{{ t('events.operations.filters.placeholderRight') }}</UiFormLabel>
					<UiFormControl>
						<VariableInput v-bind="componentField" />
					</UiFormControl>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>
		</div>
	</div>
</template>
