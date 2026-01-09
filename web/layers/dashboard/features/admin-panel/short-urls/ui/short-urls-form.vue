<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'





import {
	useAdminShortUrlsApi,
} from '~/features/admin-panel/short-urls/composables/use-admin-short-urls-api.ts'

const api = useAdminShortUrlsApi()

const formSchema = z.object({
	link: z.string().url().min(1).max(2000).trim(),
	shortId: z.string().min(1).max(10).trim().optional(),
})

const urlForm = useForm({
	validationSchema: toTypedSchema(formSchema),
	initialValues: {
		shortId: undefined,
		link: '',
	},
})

const handleSubmit = urlForm.handleSubmit(async (values) => {
	try {
		await api.createShortUrl(values.link, values.shortId)
		urlForm.resetForm()
	} catch (err) {
		console.error(err)
	}
})
</script>

<template>
	<UiCard>
		<UiCardHeader>
			<UiCardTitle>
				Create short url
			</UiCardTitle>
		</UiCardHeader>
		<UiCardContent>
			<form class="flex flex-col gap-2" @submit.prevent="handleSubmit">
				<UiFormField v-slot="{ componentField }" name="link">
					<UiFormItem>
						<UiFormLabel>Link</UiFormLabel>
						<UiFormControl>
							<UiInput type="text" v-bind="componentField" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiFormField v-slot="{ componentField }" name="shortId">
					<UiFormItem>
						<UiFormLabel>Short ID</UiFormLabel>
						<UiFormControl>
							<UiInput type="text" v-bind="componentField" />
						</UiFormControl>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>

				<UiButton type="submit" class="place-self-end mt-2">
					Create
				</UiButton>
			</form>
		</UiCardContent>
	</UiCard>
</template>
