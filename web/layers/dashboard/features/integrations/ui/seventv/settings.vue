<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { watch } from 'vue'

import * as z from 'zod'

import { useSevenTvIntegration } from '#layers/dashboard/api/integrations/seventv.ts'
import RewardsSelectorSingle from '#layers/dashboard/components/rewards-selector-single.vue'



import { toast } from 'vue-sonner'

const schema = z.object({
	rewardIdForAddEmote: z.string().nullable().default(null),
	rewardIdForRemoveEmote: z.string().nullable().default(null),
	deleteEmotesOnlyAddedByApp: z.boolean().default(true),
})

const settingsForm = useForm({
	validationSchema: toTypedSchema(schema),
	keepValuesOnUnmount: true,
	validateOnMount: false,
})

const { updater, subscription } = useSevenTvIntegration()
const { t } = useI18n()

const handleSubmit = settingsForm.handleSubmit(async (values) => {
	await updater.executeMutation({
		input: values,
	})
	toast.success(t('sharedTexts.saved'), {
		duration: 2500,
	})
})

watch(
	subscription.data,
	(data) => {
		if (data?.sevenTvData) {
			settingsForm.setValues({
				rewardIdForAddEmote: data.sevenTvData.rewardIdForAddEmote,
				rewardIdForRemoveEmote: data.sevenTvData.rewardIdForRemoveEmote,
				deleteEmotesOnlyAddedByApp: data.sevenTvData.deleteEmotesOnlyAddedByApp,
			})
		}
	},
	{ once: true }
)
</script>

<template>
	<form class="flex flex-col gap-4" @submit="handleSubmit">
		<UiFormField v-slot="{ componentField }" name="deleteEmotesOnlyAddedByApp">
			<UiFormItem class="flex flex-row items-center gap-2">
				<UiFormLabel> Delete emotes only added by app </UiFormLabel>
				<UiFormControl>
					<UiCheckbox
						:checked="componentField.modelValue"
						@update:checked="componentField['onUpdate:modelValue']"
					/>
				</UiFormControl>
				<UiFormDescription />
				<UiFormMessage />
			</UiFormItem>
		</UiFormField>

		<UiFormField v-slot="{ componentField }" name="rewardIdForAddEmote">
			<UiFormItem>
				<UiFormLabel>Select reward to listen for adding emote</UiFormLabel>
				<UiFormControl>
					<RewardsSelectorSingle
						v-model:model-value="componentField.modelValue"
						deselect
						require-input
						@update:model-value="componentField['onUpdate:modelValue']"
					/>
				</UiFormControl>
				<UiFormDescription />
				<UiFormMessage />
			</UiFormItem>
		</UiFormField>

		<UiFormField v-slot="{ componentField }" name="rewardIdForRemoveEmote">
			<UiFormItem>
				<UiFormLabel>Select reward to listen for removing emote</UiFormLabel>
				<UiFormControl>
					<RewardsSelectorSingle
						v-model:model-value="componentField.modelValue"
						deselect
						require-input
						@update:model-value="componentField['onUpdate:modelValue']"
					/>
				</UiFormControl>
				<UiFormDescription />
				<UiFormMessage />
			</UiFormItem>
		</UiFormField>

		<UiButton type="submit"> Save </UiButton>
	</form>
</template>
