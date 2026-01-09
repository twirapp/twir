<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { PlusIcon } from 'lucide-vue-next'
import { useForm } from 'vee-validate'
import { ref } from 'vue'

import { z } from 'zod'

import type { Giveaway } from '#layers/dashboard/api/giveaways.ts'

import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useGiveaways } from '~/features/giveaways/composables/giveaways-use-giveaways.ts'
import { ChannelRolePermissionEnum } from '~/gql/graphql.ts'

const { t } = useI18n()
const open = ref(false)
const { createGiveaway, viewGiveaway } = useGiveaways()

const canManageGiveaways = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageGiveaways)

// Form validation schema
const formSchema = toTypedSchema(z.object({
	keyword: z.string()
		.min(3, 'Keyword must be at least 3 characters')
		.max(100, 'Keyword must be at most 100 characters'),
}))

// Form setup
const giveawayCreateForm = useForm({
	validationSchema: formSchema,
	initialValues: {
		keyword: '',
	},
	validateOnMount: false,
})

const handleSubmit = giveawayCreateForm.handleSubmit(async (values) => {
	try {
		const result = await createGiveaway(values.keyword)
		if (result) {
			giveawayCreateForm.resetForm()
			viewGiveaway((result as Giveaway).id)
		}
	} catch (error) {
		console.error(error)
	}
})
</script>

<template>
	<UiDialog v-model:open="open">
		<UiDialogTrigger as-child>
			<UiButton size="sm" class="flex gap-2 items-center" :disabled="!canManageGiveaways">
				<PlusIcon class="size-4" />
				{{ t('giveaways.createNew') }}
			</UiButton>
		</UiDialogTrigger>

		<UiDialogContent class="sm:max-w-[425px]">
			<UiDialogHeader>
				<UiDialogTitle>{{ t('giveaways.createDialog.title') }}</UiDialogTitle>
				<UiDialogDescription>
					{{ t('giveaways.createDialog.description') }}
				</UiDialogDescription>
			</UiDialogHeader>

			<form class="space-y-4" @submit.prevent="handleSubmit">
				<UiFormField
					v-slot="{ componentField, errorMessage }"
					name="keyword"
				>
					<UiFormItem>
						<UiFormLabel>{{ t('giveaways.createDialog.keywordLabel') }}</UiFormLabel>
						<UiFormControl>
							<UiInput
								:placeholder="t('giveaways.createDialog.keywordPlaceholder')"
								v-bind="componentField"
							/>
						</UiFormControl>
						<UiFormMessage>{{ errorMessage }}</UiFormMessage>
					</UiFormItem>
				</UiFormField>

				<UiDialogFooter>
					<UiButton
						type="button"
						variant="outline"
						@click="open = false"
					>
						{{ t('giveaways.createDialog.cancel') }}
					</UiButton>
					<UiButton
						type="submit"
					>
						{{ t('giveaways.createDialog.create') }}
					</UiButton>
				</UiDialogFooter>
			</form>
		</UiDialogContent>
	</UiDialog>
</template>
