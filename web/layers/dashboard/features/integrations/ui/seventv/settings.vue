<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import { watch } from 'vue'
import { useI18n } from 'vue-i18n'
import * as z from 'zod'

import { useSevenTvIntegration } from '@/api/integrations/seventv.ts'
import RewardsSelectorSingle from '@/components/rewards-selector-single.vue'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
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
		<FormField v-slot="{ componentField }" name="deleteEmotesOnlyAddedByApp">
			<FormItem class="flex flex-row items-center gap-2">
				<FormLabel> Delete emotes only added by app </FormLabel>
				<FormControl>
					<Checkbox
						:checked="componentField.modelValue"
						@update:checked="componentField['onUpdate:modelValue']"
					/>
				</FormControl>
				<FormDescription />
				<FormMessage />
			</FormItem>
		</FormField>

		<FormField v-slot="{ componentField }" name="rewardIdForAddEmote">
			<FormItem>
				<FormLabel>Select reward to listen for adding emote</FormLabel>
				<FormControl>
					<RewardsSelectorSingle
						v-model:model-value="componentField.modelValue"
						deselect
						require-input
						@update:model-value="componentField['onUpdate:modelValue']"
					/>
				</FormControl>
				<FormDescription />
				<FormMessage />
			</FormItem>
		</FormField>

		<FormField v-slot="{ componentField }" name="rewardIdForRemoveEmote">
			<FormItem>
				<FormLabel>Select reward to listen for removing emote</FormLabel>
				<FormControl>
					<RewardsSelectorSingle
						v-model:model-value="componentField.modelValue"
						deselect
						require-input
						@update:model-value="componentField['onUpdate:modelValue']"
					/>
				</FormControl>
				<FormDescription />
				<FormMessage />
			</FormItem>
		</FormField>

		<Button type="submit"> Save </Button>
	</form>
</template>
