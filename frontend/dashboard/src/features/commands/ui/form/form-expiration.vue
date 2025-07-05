<script setup lang="ts">
import { HourglassIcon, XIcon } from 'lucide-vue-next'
import { useField } from 'vee-validate'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

import { Alert, AlertDescription } from '@/components/ui/alert'
import Button from '@/components/ui/button/Button.vue'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { DatePicker } from '@/components/ui/date-picker'
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import {
	Select,
	SelectContent,
	SelectItem,
	SelectLabel,
	SelectTrigger,
	SelectValue,
} from '@/components/ui/select'
import { useCommandEditV2 } from '@/features/commands/composables/use-command-edit-v2'
import { CommandExpiresType } from '@/gql/graphql'

const { t } = useI18n()
const { isCustom } = useCommandEditV2()

const expiresTypeOptions = computed(() => {
	return [
		{ label: t('commands.modal.expiration.actions.disable'), value: CommandExpiresType.Disable },
		{ label: t('commands.modal.expiration.actions.delete'), value: CommandExpiresType.Delete },
	]
})

const { resetField: resetAt } = useField('expiresAt')
const { resetField: resetType } = useField('expiresType')

function reset() {
	resetAt()
	resetType()
}
</script>

<template>
	<Card>
		<CardHeader class="flex flex-row place-content-center flex-wrap p-4 border-b">
			<CardTitle class="flex items-center gap-2">
				<HourglassIcon />
				{{ t('commands.modal.expiration.label') }}
			</CardTitle>
		</CardHeader>

		<CardContent v-if="isCustom" class="flex flex-col gap-4 pt-4">
			<FormField v-slot="{ componentField }" name="expiresType">
				<FormItem>
					<FormLabel>{{ t('commands.modal.expiration.actionsLabel') }}</FormLabel>
					<div class="flex flex-row gap-2">
						<FormControl>
							<Select v-bind="componentField">
								<SelectTrigger>
									<SelectValue placeholder="No expiration" />
								</SelectTrigger>
								<SelectContent>
									<SelectItem v-for="option of expiresTypeOptions" :key="option.label" :value="option.value">
										<SelectLabel>{{ option.label }}</SelectLabel>
									</SelectItem>
								</SelectContent>
							</Select>
						</FormControl>
						<Button
							variant="outline"
							type="button"
							@click="reset"
						>
							<XIcon class="size-4" />
						</Button>
					</div>
					<FormMessage />
				</FormItem>
			</FormField>

			<FormField v-slot="{ field }" name="expiresAt">
				<FormItem>
					<FormLabel>Expires at</FormLabel>
					<FormControl>
						<DatePicker
							:uid="field.name"
							auto-apply
							model-type="timestamp"
							dark
							:model-value="field.value"
							:min-date="new Date()"
							:config="{ closeOnAutoApply: true }"
							placeholder="Select date"
							@update:model-value="field['onUpdate:modelValue']"
						/>
					</FormControl>
					<FormMessage />
				</FormItem>
			</FormField>
		</CardContent>

		<CardContent v-else>
			<Alert>
				<AlertDescription>
					{{ t('commands.modal.expiration.defaultWarning') }}
				</AlertDescription>
			</Alert>
		</CardContent>
	</Card>
</template>
