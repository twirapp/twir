<script setup lang="ts">
import { Trash2 } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { Button } from '@/components/ui/button'
import {
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import VariableInput from '@/components/variable-input.vue'

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
			<Button
				type="button"
				variant="destructive"
				size="sm"
				@click="onRemove(operationIndex, filterIndex)"
			>
				<Trash2 class="h-4 w-4" />
			</Button>
		</div>

		<div class="space-y-4">
			<FormField
				v-slot="{ componentField }"
				:name="`operations.${operationIndex}.filters.${filterIndex}.left`"
			>
				<FormItem>
					<FormLabel>{{ t('events.operations.filters.placeholderLeft') }}</FormLabel>
					<FormControl>
						<VariableInput v-bind="componentField" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField
				v-slot="{ componentField }"
				:name="`operations.${operationIndex}.filters.${filterIndex}.type`"
			>
				<FormItem>
					<FormControl>
						<Select
							v-bind="componentField"
							:placeholder="t('events.selectFilterType')"
						>
							<SelectTrigger>
								<SelectValue :placeholder="t('events.selectFilterType')" />
							</SelectTrigger>
							<SelectContent>
								<SelectItem value="EQUALS">
									=
								</SelectItem>
								<SelectItem value="NOT_EQUALS">
									!=
								</SelectItem>
								<SelectItem value="CONTAINS">
									contains
								</SelectItem>
								<SelectItem value="NOT_CONTAINS">
									not contains
								</SelectItem>
								<SelectItem value="STARTS_WITH">
									starts with
								</SelectItem>
								<SelectItem value="ENDS_WITH">
									ends with
								</SelectItem>
								<SelectItem value="GREATER_THAN">
									>
								</SelectItem>
								<SelectItem value="LESS_THAN">
									{{ '<' }}
								</SelectItem>
								<SelectItem value="GREATER_THAN_OR_EQUALS">
									>=
								</SelectItem>
								<SelectItem value="LESS_THAN_OR_EQUALS">
									{{ '<=' }}
								</SelectItem>
							</SelectContent>
						</Select>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField
				v-slot="{ componentField }"
				:name="`operations.${operationIndex}.filters.${filterIndex}.right`"
			>
				<FormItem>
					<FormLabel>{{ t('events.operations.filters.placeholderRight') }}</FormLabel>
					<FormControl>
						<VariableInput v-bind="componentField" />
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>
		</div>
	</div>
</template>
