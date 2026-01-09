<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'

import * as z from 'zod'

import { useMutationEventSubSubscribe } from '#layers/dashboard/api/admin/actions'





import { toast } from 'vue-sonner'

const { t } = useI18n()

const mutationEventSubSubscribe = useMutationEventSubSubscribe()

const formSchema = toTypedSchema(
	z.object({
		type: z.string(),
		version: z.string(),
	})
)

const { handleSubmit } = useForm({
	validationSchema: formSchema,
	initialValues: {
		version: '1',
	},
})

const onSubmit = handleSubmit(async (values) => {
	const result = await mutationEventSubSubscribe.executeMutation({
		opts: values,
	})

	if (result.error) {
		toast.error(result.error.message, {
			duration: 2500,
		})
	}
})
</script>

<template>
	<h4 class="scroll-m-20 text-xl font-semibold tracking-tight">
		{{ t('adminPanel.adminActions.eventsub.title') }}
	</h4>

	<UiCard>
		<form @submit.prevent="onSubmit">
			<UiCardContent class="p-4">
				<div class="grid items-center w-full gap-4">
					<UiFormField v-slot="{ componentField }" name="version">
						<UiFormItem>
							<UiLabel for="version">
								{{ t('adminPanel.adminActions.eventsub.version') }}
							</UiLabel>
							<UiFormControl>
								<UiInput v-bind="componentField" />
							</UiFormControl>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>

					<UiFormField v-slot="{ componentField }" name="type">
						<UiFormItem>
							<UiLabel for="type">
								{{ t('adminPanel.adminActions.eventsub.type') }}
							</UiLabel>
							<UiFormControl>
								<UiInput v-bind="componentField" />
							</UiFormControl>
							<UiFormDescription>
								You can find all available subscription types in
								<a
									class="underline"
									href="https://dev.twitch.tv/docs/eventsub/eventsub-subscription-types/"
									target="_blank"
									>subscription types</a
								>.
							</UiFormDescription>
							<UiFormMessage />
						</UiFormItem>
					</UiFormField>
				</div>
			</UiCardContent>

			<UiCardFooter class="flex justify-end p-4">
				<UiButton>
					{{ t('sharedButtons.send') }}
				</UiButton>
			</UiCardFooter>
		</form>
	</UiCard>
</template>
