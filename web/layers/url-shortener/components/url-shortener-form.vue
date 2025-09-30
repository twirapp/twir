<script setup lang="ts">
import { toTypedSchema } from '@vee-validate/zod'
import { useForm } from 'vee-validate'
import * as z from 'zod'

import UrlShortenerShortenUrl from './url-shortener-shorten-url.vue'
import { useUrlShortener } from '../composables/use-url-shortener'
import type { LinkOutputDto } from '@twir/api/openapi'

const schema = z.object({
	url: z.url('Please enter a valid URL'),
	customAlias: z
		.string()
		.min(3, 'Custom alias must be at least 3 characters')
		.max(30, 'Custom alias must be at most 30 characters')
		.regex(/^[a-zA-Z0-9]*$/, 'Custom alias can only contain letters and numbers')
		.optional()
		.or(z.literal('')),
})

const shortenerForm = useForm({
	// @ts-ignore
	validationSchema: toTypedSchema(schema),
})

const api = useUrlShortener()

const currentUrl = ref<LinkOutputDto>()
const currentError = ref<string>()

const onSubmit = shortenerForm.handleSubmit(async (values) => {
	currentUrl.value = undefined
	currentError.value = undefined

	const { data, error } = await api.shortUrl({
		url: values.url,
		alias: values.customAlias || undefined,
	})

	if (error) {
		currentError.value = error.value.message
		return
	}

	currentUrl.value = data
})
</script>

<template>
	<UiCard class="w-full lg:w-[40%] bg-black/60">
		<UiCardHeader class="text-2xl font-bold text-center"> URL Shortener </UiCardHeader>
		<UiCardContent>
			<form class="flex flex-col gap-2" @submit.prevent="onSubmit">
				<UiFormField v-slot="{ componentField }" name="url">
					<UiFormItem>
						<UiFormLabel>URL</UiFormLabel>
						<UiFormControl>
							<UiInput
								type="text"
								placeholder="https://twitch.tv/twirdev"
								v-bind="componentField"
							/>
						</UiFormControl>
						<UiFormDescription> URL you want to short. </UiFormDescription>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
				<UiFormField v-slot="{ componentField }" name="customAlias">
					<UiFormItem>
						<UiFormLabel>
							Custom Alias <span class="text-xs font-normal">(Optional)</span>
						</UiFormLabel>
						<UiFormControl>
							<UiInput type="text" placeholder="stream" v-bind="componentField" />
						</UiFormControl>
						<UiFormDescription>
							Custom alias for your shortened URL. Must be between 3 and 30 characters.
						</UiFormDescription>
						<UiFormMessage />
					</UiFormItem>
				</UiFormField>
				<UiButton type="submit" class="mt-4">Short URL</UiButton>
			</form>
		</UiCardContent>
		<UiCardFooter>
			<div
				class="block p-4 border-border border-2 rounded-md bg-red-900 w-full underline"
				v-if="currentError"
			>
				{{ currentError }}
			</div>
			<UrlShortenerShortenUrl v-else-if="currentUrl" :url="currentUrl" />
		</UiCardFooter>
	</UiCard>
</template>
