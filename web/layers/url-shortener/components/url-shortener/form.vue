<script setup lang="ts">
import { useForm } from 'vee-validate'
import * as z from 'zod'
import { storeToRefs } from 'pinia'

import { useUrlShortener } from '../../composables/use-url-shortener'

import type { LinkOutputDto } from '@twir/api/openapi'

const twirShortenerUrl = useRequestURL()

const schema = z.object({
	url: z
		.url({
			protocol: /^https?$/,
			hostname: z.regexes.domain,
		})
		.default(''),
	customAlias: z
		.string()
		.min(3, 'Custom alias must be at least 3 characters')
		.max(30, 'Custom alias must be at most 30 characters')
		.regex(/^[a-zA-Z0-9]*$/, 'Custom alias can only contain letters and numbers')
		.optional()
		.or(z.literal('')),
	useCustomDomain: z.boolean().default(false),
})

const form = useForm({
	validationSchema: schema,
	initialValues: {
		useCustomDomain: false,
	},
})

const api = useUrlShortener()
const { customDomain } = storeToRefs(api)
const currentUrl = ref<LinkOutputDto>()
const currentError = ref<string>()
const { setFieldValue, values } = form
const hasVerifiedCustomDomain = computed(
	() => Boolean(customDomain.value?.domain && customDomain.value?.verified)
)
const hasManualDomainChoice = ref(false)
const useCustomDomainPreference = computed(() => {
	if (!hasVerifiedCustomDomain.value) return false
	return Boolean(values.useCustomDomain)
})
const customDomainPrefix = computed(() => {
	if (useCustomDomainPreference.value && customDomain.value?.domain) {
		return `https://${customDomain.value.domain}/`
	}

	return `${twirShortenerUrl.origin}/s/`
})

function handleCustomDomainToggle(
	nextValue: boolean,
	handleChange: (value: boolean) => void
) {
	hasManualDomainChoice.value = true
	handleChange(nextValue)
}

watch(
	hasVerifiedCustomDomain,
	(hasCustomDomain) => {
		if (!hasCustomDomain) {
			hasManualDomainChoice.value = false
			setFieldValue('useCustomDomain', false)
			return
		}
		if (!hasManualDomainChoice.value) {
			setFieldValue('useCustomDomain', true)
		}
	},
	{ immediate: true }
)

const onSubmit = form.handleSubmit(async (values) => {
	currentUrl.value = undefined
	currentError.value = undefined

	const { data, error } = await api.shortUrl({
		url: values.url,
		alias: values.customAlias || undefined,
		useCustomDomain: useCustomDomainPreference.value,
	})

	if (error) {
		currentError.value = error
		return
	}

	if (data?.data) {
		currentUrl.value = data.data
		form.resetForm({
			values: {
				url: '',
				customAlias: '',
				useCustomDomain: values.useCustomDomain ?? false,
			},
		})
	}
})
</script>

<template>
	<div
		class="flex flex-col items-start border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,9%)] h-fit w-full max-w-xl rounded-2xl p-3 shadow-[0px_0px_30px_hsl(240,11%,6%)]"
	>
		<form class="w-full" @submit="onSubmit">
			<!-- Input -->
			<UiFormField v-slot="{ componentField }" name="url">
				<UiFormItem>
					<UiFormControl>
						<div
							class="flex flex-col items-start sm:flex-row sm:items-center rounded-xl p-2 border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,15%)] w-full focus-within:ring ring-[hsl(240,11%,30%)] focus-within:border-[hsl(240,11%,30%)] focus-within:ring-[hsl(240,11%,30%)] transition-all"
						>
							<label for="url" class="flex items-center w-full gap-2.5 overflow-hidden">
								<Icon name="lucide:link" class="flex-none w-4 h-4 ml-2 mr-1 cursor-pointer" />
								<input
									id="url"
									type="text"
									v-bind="componentField"
									class="flex-1 bg-transparent font-medium focus-visible:outline-none placeholder-[hsl(240,11%,40%)]"
									placeholder="https://twitch.tv/twirdev"
								/>
							</label>
							<UiButton
								type="submit"
								variant="outline"
								class="h-fit w-full sm:w-auto mt-2 sm:mt-0 py-1 sm:py-1.5 sm:px-3 rounded-lg font-semibold border border-[hsl(240,11%,30%)] hover:border-[hsl(240,11%,45%)] bg-[hsl(240,11%,25%)] hover:bg-[hsl(240,11%,35%)]"
							>
								Shorten
							</UiButton>
						</div>
					</UiFormControl>
					<UiFormMessage class="px-2 pt-1 text-xs" />
				</UiFormItem>
			</UiFormField>

			<!-- Optional -->
			<UiAccordion type="single" collapsible class="w-full">
				<UiAccordionItem value="item-1">
					<UiAccordionTrigger
						class="px-2 text-sm text-[hsl(240,11%,70%)] hover:no-underline hover:text-white"
					>
						<span class="flex gap-2 items-center"> Optional </span>
					</UiAccordionTrigger>
					<UiAccordionContent class="accordion-content">
						<div v-if="hasVerifiedCustomDomain" class="px-2 pb-3">
							<UiFormField v-slot="{ value, handleChange }" name="useCustomDomain">
								<UiFormItem class="flex items-center justify-between gap-3">
									<div class="space-y-1">
										<UiFormLabel class="text-sm">Use custom domain</UiFormLabel>
										<UiFormDescription class="text-xs text-[hsl(240,11%,55%)]">
											Default: {{ customDomain?.domain }}
										</UiFormDescription>
									</div>
									<UiFormControl>
										<UiSwitch
											:model-value="value"
											@update:model-value="
												(nextValue: boolean) => handleCustomDomainToggle(nextValue, handleChange)
											"
										/>
									</UiFormControl>
									<UiFormMessage />
								</UiFormItem>
							</UiFormField>
						</div>
						<UiFormField v-slot="{ componentField }" name="customAlias">
							<UiFormItem>
								<div class="flex flex-col gap-2 px-2">
									<UiFormLabel
										for="alias"
										class="flex items-center gap-3 font-medium cursor-pointer"
									>
										Custom alias
									</UiFormLabel>
									<UiFormControl>
										<div
											class="flex items-center rounded-xl p-1.5 py-3 border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,15%)] w-full"
										>
											<label for="alias" class="flex items-center w-full">
												<span
													class="ml-2 mr-1 font-semibold"
													:class="{
														'text-yellow-300/80': customDomain && !customDomain.verified,
													}"
												>
													{{ customDomainPrefix }}
												</span>
												<input
													id="alias"
													type="text"
													v-bind="componentField"
													class="flex-1 font-extrabold bg-transparent focus-visible:outline-none placeholder-[hsl(240,11%,40%)]"
													placeholder="stream"
												/>
											</label>
										</div>
									</UiFormControl>
									<UiFormDescription class="text-xs text-[hsl(240,11%,50%)]">
										Custom alias for your shortened URL. Must be between 3 and 30 characters.
									</UiFormDescription>
									<p
										v-if="customDomain && !customDomain.verified"
										class="text-xs text-yellow-300/80"
									>
										Verify your custom domain to use it for new links. Until then, links use
										{{ twirShortenerUrl.origin + '/s/' }}.
									</p>
									<UiFormMessage class="text-xs" />
								</div>
							</UiFormItem>
						</UiFormField>
					</UiAccordionContent>
				</UiAccordionItem>
			</UiAccordion>

			<!-- Error Display -->
			<div v-if="currentError" class="px-2 pt-2 text-sm text-red-500">
				{{ currentError }}
			</div>
		</form>
	</div>
</template>
