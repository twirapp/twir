<script setup lang="ts">
import { Trash2 } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'

import { Button } from '@/components/ui/button'
import {
	FormControl,
	FormDescription,
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

const props = defineProps<{
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
				{{ t('events.filter') }} {{ filterIndex + 1 }}
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
				:name="`operations.${operationIndex}.filters.${filterIndex}.type`"
			>
				<FormItem>
					<FormLabel>{{ t('events.filterType') }}</FormLabel>
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
									{{ t('events.filterEquals') }}
								</SelectItem>
								<SelectItem value="NOT_EQUALS">
									{{ t('events.filterNotEquals') }}
								</SelectItem>
								<SelectItem value="CONTAINS">
									{{ t('events.filterContains') }}
								</SelectItem>
								<SelectItem value="NOT_CONTAINS">
									{{ t('events.filterNotContains') }}
								</SelectItem>
								<SelectItem value="STARTS_WITH">
									{{ t('events.filterStartsWith') }}
								</SelectItem>
								<SelectItem value="ENDS_WITH">
									{{ t('events.filterEndsWith') }}
								</SelectItem>
								<SelectItem value="GREATER_THAN">
									{{ t('events.filterGreaterThan') }}
								</SelectItem>
								<SelectItem value="LESS_THAN">
									{{ t('events.filterLessThan') }}
								</SelectItem>
								<SelectItem value="GREATER_THAN_OR_EQUALS">
									{{ t('events.filterGreaterThanOrEquals') }}
								</SelectItem>
								<SelectItem value="LESS_THAN_OR_EQUALS">
									{{ t('events.filterLessThanOrEquals') }}
								</SelectItem>
							</SelectContent>
						</Select>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField
				v-slot="{ componentField }"
				:name="`operations.${operationIndex}.filters.${filterIndex}.left`"
			>
				<FormItem>
					<FormLabel>{{ t('events.filterLeft') }}</FormLabel>
					<FormControl>
						<VariableInput v-bind="componentField" />
					</FormControl>
					<FormDescription>{{ t('events.filterLeftDescription') }}</FormDescription>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField
				v-slot="{ componentField }"
				:name="`operations.${operationIndex}.filters.${filterIndex}.right`"
			>
				<FormItem>
					<FormLabel>{{ t('events.filterRight') }}</FormLabel>
					<FormControl>
						<VariableInput v-bind="componentField" />
					</FormControl>
					<FormDescription>{{ t('events.filterRightDescription') }}</FormDescription>
					<FormMessage />
				</FormItem>
			</FormField>
		</div>
	</div>
</template>
