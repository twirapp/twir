<script setup lang="ts">
import { RadioGroupItem, RadioGroupRoot } from 'reka-ui'
import { useForm } from 'vee-validate'
import * as z from 'zod'

import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from '@/components/ui/card'
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from '@/components/ui/form'
import { Skeleton } from '@/components/ui/skeleton'
import { parseMcpConsentAttempt } from '~/utils/mcp-consent.js'

import { createMcpConsentApi, type McpConsent } from './api.js'

type ConsentScreenState = 'loading' | 'ready' | 'expired' | 'permission' | 'network'

const { t } = useI18n()
const route = useRoute()
const mcpConsentApi = createMcpConsentApi($fetch)
const attempt = computed(() => parseMcpConsentAttempt(route.query.attempt))
const consent = ref<McpConsent | null>(null)
const screenState = ref<ConsentScreenState>('loading')
const isSubmitting = ref(false)
const hasCompleted = ref(false)

const formSchema = z.object({
	accessLevel: z.enum(['read', 'write']),
})

const { handleSubmit, values } = useForm({
	validationSchema: formSchema,
	initialValues: { accessLevel: 'read' },
})

const selectedDashboardName = computed(() => consent.value?.channel_id ?? '')

const isActionDisabled = computed(() => isSubmitting.value || hasCompleted.value)

function setScreenState(result: Awaited<ReturnType<typeof mcpConsentApi.getMcpConsent>>): void {
	switch (result.kind) {
		case 'success':
			consent.value = result.data
			screenState.value = 'ready'
			return
		case 'expired':
		case 'permission':
		case 'network':
			screenState.value = result.kind
			return
		default:
			return result satisfies never
	}
}

async function loadConsent(): Promise<void> {
	const currentAttempt = attempt.value
	if (currentAttempt === null) {
		screenState.value = 'expired'
		return
	}

	screenState.value = 'loading'
	setScreenState(await mcpConsentApi.getMcpConsent(currentAttempt))
}

async function submitDecision(
	decision: 'approve' | 'deny',
	accessLevel?: 'read' | 'write',
): Promise<void> {
	const currentConsent = consent.value
	const currentAttempt = attempt.value
	if (currentConsent === null || currentAttempt === null || isActionDisabled.value) return

	isSubmitting.value = true
	const result = await mcpConsentApi.submitMcpConsent(
		decision === 'approve'
			? {
					attempt: currentAttempt,
					csrf_token: currentConsent.csrf_token,
					channel_id: currentConsent.channel_id,
					decision,
					access_level: accessLevel ?? 'read',
				}
			: {
					attempt: currentAttempt,
					csrf_token: currentConsent.csrf_token,
					channel_id: currentConsent.channel_id,
					decision,
				},
	)

	switch (result.kind) {
		case 'success':
			hasCompleted.value = true
			window.location.replace(result.data.redirectTo)
			return
		case 'expired':
		case 'permission':
		case 'network':
			screenState.value = result.kind
			isSubmitting.value = false
			return
		default:
			return result satisfies never
	}
}

const approve = handleSubmit(async (formValues) => {
	await submitDecision('approve', formValues.accessLevel)
})

async function deny(): Promise<void> {
	await submitDecision('deny')
}

onMounted(() => {
	void loadConsent()
})
</script>

<template>
	<PageLayout>
		<template #title>{{ t('mcpConsent.title') }}</template>

		<template #title-footer>
			<p class="text-muted-foreground">{{ t('mcpConsent.description') }}</p>
		</template>

		<template #content>
			<div class="mx-auto flex w-full max-w-2xl flex-col gap-6">
				<Card v-if="screenState === 'loading'">
					<CardHeader>
						<Skeleton class="h-6 w-48" />
						<Skeleton class="h-4 w-full" />
					</CardHeader>
					<CardContent class="flex flex-col gap-4">
						<Skeleton class="h-20 w-full" />
						<Skeleton class="h-10 w-full" />
					</CardContent>
				</Card>

				<Alert v-else-if="screenState === 'expired'" variant="destructive">
					<Icon name="lucide:clock-alert" />
					<AlertTitle>{{ t('mcpConsent.errors.expired.title') }}</AlertTitle>
					<AlertDescription>{{ t('mcpConsent.errors.expired.description') }}</AlertDescription>
				</Alert>

				<Alert v-else-if="screenState === 'permission'" variant="destructive">
					<Icon name="lucide:shield-alert" />
					<AlertTitle>{{ t('mcpConsent.errors.permission.title') }}</AlertTitle>
					<AlertDescription>{{ t('mcpConsent.errors.permission.description') }}</AlertDescription>
				</Alert>

				<Alert v-else-if="screenState === 'network'" variant="destructive">
					<Icon name="lucide:triangle-alert" />
					<AlertTitle>{{ t('mcpConsent.errors.network.title') }}</AlertTitle>
					<AlertDescription class="flex flex-col gap-3">
						<span>{{ t('mcpConsent.errors.network.description') }}</span>
						<Button class="w-fit" type="button" variant="outline" @click="loadConsent">
							{{ t('mcpConsent.retry') }}
						</Button>
					</AlertDescription>
				</Alert>

				<form v-else-if="consent" class="flex flex-col" @submit="approve">
					<Card>
						<CardHeader>
							<CardTitle>{{
								t('mcpConsent.requestTitle', { client: consent.client.name })
							}}</CardTitle>
							<CardDescription>{{ t('mcpConsent.requestDescription') }}</CardDescription>
						</CardHeader>
						<CardContent class="flex flex-col gap-6">
							<dl class="grid gap-4 text-sm sm:grid-cols-2">
								<div class="flex flex-col gap-1">
									<dt class="text-muted-foreground">{{ t('mcpConsent.client') }}</dt>
									<dd class="font-medium">{{ consent.client.name }}</dd>
								</div>
								<div class="flex flex-col gap-1">
									<dt class="text-muted-foreground">{{ t('mcpConsent.dashboard') }}</dt>
									<dd class="font-medium">{{ selectedDashboardName }}</dd>
								</div>
								<div v-if="consent.client.uri" class="flex flex-col gap-1 sm:col-span-2">
									<dt class="text-muted-foreground">{{ t('mcpConsent.clientUri') }}</dt>
									<dd class="break-all font-medium">{{ consent.client.uri }}</dd>
								</div>
								<div class="flex flex-col gap-1 sm:col-span-2">
									<dt class="text-muted-foreground">{{ t('mcpConsent.requestedScopes') }}</dt>
									<dd class="font-medium">{{ consent.requested_scopes.join(', ') }}</dd>
								</div>
							</dl>

							<FormField v-slot="{ value, handleChange }" name="accessLevel">
								<FormItem>
									<FormLabel>{{ t('mcpConsent.accessLevel') }}</FormLabel>
									<FormDescription>{{ t('mcpConsent.readDescription') }}</FormDescription>
									<FormControl>
										<RadioGroupRoot
											:model-value="value"
											:disabled="isActionDisabled"
											class="grid gap-3"
											:aria-label="t('mcpConsent.accessLevel')"
											@update:model-value="handleChange"
										>
											<div class="flex items-start gap-3 rounded-lg border p-4">
												<RadioGroupItem id="mcp-consent-read" class="mt-0.5" value="read" />
												<label class="flex cursor-pointer flex-col gap-1" for="mcp-consent-read">
													<span class="font-medium">{{ t('mcpConsent.readTitle') }}</span>
													<span class="text-muted-foreground text-sm">{{
														t('mcpConsent.readDescription')
													}}</span>
												</label>
											</div>

											<div
												v-if="consent.access_levels.includes('write')"
												class="flex items-start gap-3 rounded-lg border p-4"
											>
												<RadioGroupItem id="mcp-consent-write" class="mt-0.5" value="write" />
												<label class="flex cursor-pointer flex-col gap-1" for="mcp-consent-write">
													<span class="font-medium">{{ t('mcpConsent.writeTitle') }}</span>
													<span class="text-muted-foreground text-sm">{{
														t('mcpConsent.writeDescription')
													}}</span>
												</label>
											</div>
										</RadioGroupRoot>
									</FormControl>
									<FormMessage />
								</FormItem>
							</FormField>

							<Alert v-if="values.accessLevel === 'write'" variant="destructive">
								<Icon name="lucide:shield-alert" />
								<AlertTitle>{{ t('mcpConsent.writeWarning.title') }}</AlertTitle>
								<AlertDescription>{{ t('mcpConsent.writeWarning.description') }}</AlertDescription>
							</Alert>
						</CardContent>
						<CardFooter class="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
							<Button :disabled="isActionDisabled" type="button" variant="outline" @click="deny">
								{{ t('mcpConsent.deny') }}
							</Button>
							<Button :disabled="isActionDisabled" type="submit">
								<Icon
									v-if="isSubmitting"
									name="lucide:loader-circle"
									class="animate-spin"
									data-icon="inline-start"
								/>
								{{ t('mcpConsent.approve') }}
							</Button>
						</CardFooter>
					</Card>
				</form>
			</div>
		</template>
	</PageLayout>
</template>
