<script setup lang="ts">
import { HourglassIcon, XIcon } from 'lucide-vue-next'
import { useField } from 'vee-validate'
import { computed } from 'vue'








import { useCommandEditV2 } from '~/features/commands/composables/use-command-edit-v2'
import { CommandExpiresType } from '~/gql/graphql'

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
	<UiCard>
		<UiCardHeader class="flex flex-row place-content-center flex-wrap">
			<UiCardTitle class="flex items-center gap-2">
				<HourglassIcon />
				{{ t('commands.modal.expiration.label') }}
			</UiCardTitle>
		</UiCardHeader>

		<UiCardContent v-if="isCustom" class="flex flex-col gap-4 pt-4">
			<UiFormField v-slot="{ componentField }" name="expiresType">
				<UiFormItem>
					<UiFormLabel>{{ t('commands.modal.expiration.actionsLabel') }}</UiFormLabel>
					<div class="flex flex-row gap-2">
						<UiFormControl>
							<UiSelect v-bind="componentField">
								<UiSelectTrigger>
									<UiSelectValue placeholder="No expiration" />
								</UiSelectTrigger>
								<UiSelectContent>
									<UiSelectItem
										v-for="option of expiresTypeOptions"
										:key="option.label"
										:value="option.value"
									>
										<UiSelectLabel>{{ option.label }}</UiSelectLabel>
									</UiSelectItem>
								</UiSelectContent>
							</UiSelect>
						</UiFormControl>
						<UiButton variant="outline" type="button" @click="reset">
							<XIcon class="size-4" />
						</UiButton>
					</div>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>

			<UiFormField v-slot="{ field }" name="expiresAt">
				<UiFormItem>
					<UiFormLabel>Expires at</UiFormLabel>
					<UiFormControl>
						<UiDatePicker
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
					</UiFormControl>
					<UiFormMessage />
				</UiFormItem>
			</UiFormField>
		</UiCardContent>

		<UiCardContent v-else>
			<UiAlert>
				<UiAlertDescription>
					{{ t('commands.modal.expiration.defaultWarning') }}
				</UiAlertDescription>
			</UiAlert>
		</UiCardContent>
	</UiCard>
</template>
