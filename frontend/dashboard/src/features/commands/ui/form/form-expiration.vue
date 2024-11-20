<script setup lang="ts">
import { XIcon } from 'lucide-vue-next'
import { useField } from 'vee-validate'
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

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
const { command } = useCommandEditV2()

const expiresTypeOptions = computed(() => {
	return [
		{ label: t('commands.modal.expiration.actions.disable'), value: CommandExpiresType.Disable },
		{ label: t('commands.modal.expiration.actions.delete'), value: CommandExpiresType.Delete },
	]
})

const { setValue: setExpiresValue } = useField('expiresAt')
const { setValue: setTypeValue } = useField('expiresType')

function reset() {
	setExpiresValue(null)
	setTypeValue(null)
}
</script>

<template>
	<Card>
		<CardHeader class="flex flex-row justify-between flex-wrap">
			<CardTitle>{{ t('commands.modal.expiration.label') }}</CardTitle>
		</CardHeader>

		<CardContent class="flex flex-col gap-4">
			<FormField v-slot="{ componentField }" name="expiresType">
				<FormItem>
					<FormLabel>{{ t('commands.modal.expiration.actionsLabel') }}</FormLabel>
					<div class="flex flex-row gap-2">
						<FormControl>
							<Select v-bind="componentField" :disabled="command?.default">
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
							:disabled="command?.default"
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
							:disabled="command?.default"
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
	</Card>
</template>
