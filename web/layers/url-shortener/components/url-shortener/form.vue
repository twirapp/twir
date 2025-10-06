<script setup lang="ts">
import { useForm } from 'vee-validate'
import * as z from 'zod'

import { useUrlShortener } from '../../composables/use-url-shortener'

import type { LinkOutputDto } from '@twir/api/openapi'

const twirShortenerUrl = useRequestURL()

const schema = z.object({
	url: z.url({
		protocol: /^https?$/,
		hostname: z.regexes.domain,
	}),
	customAlias: z
		.string()
		.min(3, 'Custom alias must be at least 3 characters')
		.max(30, 'Custom alias must be at most 30 characters')
		.regex(/^[a-zA-Z0-9]*$/, 'Custom alias can only contain letters and numbers')
		.optional()
		.or(z.literal('')),
})

const form = useForm({
	validationSchema: schema,
})

const api = useUrlShortener()

const STORAGE_KEY = 'shortener_history'
const MAX_HISTORY = 3

const recentUrls = ref<LinkOutputDto[]>([])

onMounted(() => {
	if (process.client) {
		const stored = localStorage.getItem(STORAGE_KEY)
		if (stored) {
			try {
				recentUrls.value = JSON.parse(stored)
			} catch (e) {
				console.error('Failed to parse localStorage:', e)
				recentUrls.value = []
			}
		}
	}
})

function saveToHistory(url: LinkOutputDto) {
	recentUrls.value = [url, ...recentUrls.value]

	recentUrls.value = recentUrls.value.filter(
		(item, index, self) => index === self.findIndex((t) => t.id === item.id)
	)

	if (recentUrls.value.length > MAX_HISTORY) {
		recentUrls.value = recentUrls.value.slice(0, MAX_HISTORY)
	}

	if (process.client) {
		localStorage.setItem(STORAGE_KEY, JSON.stringify(recentUrls.value))
	}
}

const currentUrl = ref<LinkOutputDto>()
const currentError = ref<string>()

const onSubmit = form.handleSubmit(async (values) => {
	currentUrl.value = undefined
	currentError.value = undefined

	const { data, error } = await api.shortUrl({
		url: values.url,
		alias: values.customAlias || undefined,
	})

	if (error) {
		currentError.value = error
		return
	}

	if (data) {
		currentUrl.value = data
		saveToHistory(data)
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
							class="flex items-center rounded-xl p-1.5 border border-[hsl(240,11%,18%)] bg-[hsl(240,11%,15%)] w-full"
						>
							<label for="url" class="flex items-center w-full gap-2.5">
								<Icon name="lucide:link" class="size-4 ml-2 mr-1 cursor-pointer" />
								<input
									id="url"
									type="text"
									v-bind="componentField"
									class="flex-1 bg-transparent font-medium focus-visible:outline-none border-l border-[hsl(240,11%,40%)] pl-4 placeholder-[hsl(240,11%,40%)]"
									placeholder="https://twitch.tv/twirdev"
								/>
								<UiButton
									type="submit"
									variant="outline"
									class="h-fit px-3 rounded-lg font-semibold border border-[hsl(240,11%,25%)] hover:border-[hsl(240,11%,40%)] bg-[hsl(240,11%,20%)] hover:bg-[hsl(240,11%,30%)]"
								>
									Shorten
								</UiButton>
							</label>
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
												<span class="ml-2 mr-1 font-semibold">
													{{ twirShortenerUrl.origin + '/s/' }}
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

	<!-- Recent URLs -->
	<div class="flex flex-col w-full max-w-xl gap-3 mt-4">
		<UrlShortenerCard v-if="recentUrls.length === 0" />
		<TransitionGroup name="list" tag="div" class="flex flex-col gap-3">
			<div v-for="url in recentUrls" :key="url.id">
				<UrlShortenerCard :url="url" />
			</div>
		</TransitionGroup>
	</div>
</template>

<style>
.list-enter-active,
.list-leave-active {
	transition: all 0.3s ease;
}

.list-enter-from {
	opacity: 0;
	transform: translateY(-20px);
}

.list-leave-to {
	opacity: 0;
	transform: translateY(20px);
}

.list-move {
	transition: transform 0.3s ease;
}
</style>
