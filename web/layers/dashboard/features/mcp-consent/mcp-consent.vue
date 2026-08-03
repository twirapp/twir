<script setup lang="ts">
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
import { parseMcpConsentAttempt } from '~/utils/mcp-consent.js'
import {
	type McpApprovedScopes,
	type McpScopeGroup,
	buildMcpScopeToken,
	mcpApprovedScopesSchema,
} from '~/utils/mcp-scopes.js'

import { type McpConsent, type McpConsentDecision, createMcpConsentApi } from './api.js'
import McpConsentState from './mcp-consent-state.vue'
import ScopeGroupToggle from './scope-group-toggle.vue'

type ConsentScreenState = 'loading' | 'ready' | 'expired' | 'permission' | 'network'

interface GroupSelection {
	readonly read: boolean
	readonly edit: boolean
}

const { t } = useI18n()
const route = useRoute()
const mcpConsentApi = createMcpConsentApi($fetch)
const attempt = computed(() => parseMcpConsentAttempt(route.query.attempt))
const consent = ref<McpConsent | null>(null)
const screenState = ref<ConsentScreenState>('loading')
const selection = ref<Partial<Record<McpScopeGroup, GroupSelection>>>({})
const showSelectionError = ref(false)
const isSubmitting = ref(false)
const hasCompleted = ref(false)

const selectedDashboardName = computed(() => consent.value?.channel_id ?? '')

const isActionDisabled = computed(() => isSubmitting.value || hasCompleted.value)

const hasEditSelected = computed(() =>
	Object.values(selection.value).some((groupSelection) => groupSelection?.edit === true),
)

function setScreenState(result: Awaited<ReturnType<typeof mcpConsentApi.getMcpConsent>>): void {
	switch (result.kind) {
		case 'success':
			consent.value = result.data
			selection.value = Object.fromEntries(
				result.data.requested_scopes.map((scope) => [
					scope.group,
					{ read: true, edit: false },
				]),
			)
			showSelectionError.value = false
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

function setRead(group: McpScopeGroup, value: boolean): void {
	selection.value[group] = value
		? { read: true, edit: selection.value[group]?.edit === true }
		: { read: false, edit: false }
	showSelectionError.value = false
}

function setEdit(group: McpScopeGroup, value: boolean): void {
	selection.value[group] = value
		? { read: true, edit: true }
		: { read: selection.value[group]?.read === true, edit: false }
	showSelectionError.value = false
}

function buildApprovedScopes(currentConsent: McpConsent): McpApprovedScopes | null {
	const tokens = currentConsent.requested_scopes.flatMap((scope) => {
		const groupSelection = selection.value[scope.group]
		if (groupSelection?.read !== true) return []

		const readToken = buildMcpScopeToken(scope.group, 'read')
		if (groupSelection.edit && scope.actions.includes('edit')) {
			return [readToken, buildMcpScopeToken(scope.group, 'edit')]
		}

		return [readToken]
	})

	const result = mcpApprovedScopesSchema.safeParse(tokens)
	return result.success ? result.data : null
}

async function submitDecision(decision: McpConsentDecision): Promise<void> {
	if (consent.value === null || isActionDisabled.value) return

	isSubmitting.value = true
	const result = await mcpConsentApi.submitMcpConsent(decision)

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

async function approve(): Promise<void> {
	const currentConsent = consent.value
	const currentAttempt = attempt.value
	if (currentConsent === null || currentAttempt === null || isActionDisabled.value) return

	const approvedScopes = buildApprovedScopes(currentConsent)
	if (approvedScopes === null) {
		showSelectionError.value = true
		return
	}

	await submitDecision({
		attempt: currentAttempt,
		csrf_token: currentConsent.csrf_token,
		channel_id: currentConsent.channel_id,
		decision: 'approve',
		approved_scopes: approvedScopes,
	})
}

async function deny(): Promise<void> {
	const currentConsent = consent.value
	const currentAttempt = attempt.value
	if (currentConsent === null || currentAttempt === null || isActionDisabled.value) return

	await submitDecision({
		attempt: currentAttempt,
		csrf_token: currentConsent.csrf_token,
		channel_id: currentConsent.channel_id,
		decision: 'deny',
	})
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
				<McpConsentState
					v-if="screenState !== 'ready'"
					:state="screenState"
					@retry="loadConsent"
				/>

				<form v-else-if="consent" class="flex flex-col" @submit.prevent="approve">
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
							</dl>

							<fieldset class="flex flex-col gap-3">
								<legend class="sr-only">{{ t('mcpConsent.scopes.title') }}</legend>
								<div class="flex flex-col gap-1">
									<h3 class="text-sm font-medium">{{ t('mcpConsent.scopes.title') }}</h3>
									<p class="text-muted-foreground text-sm">
										{{ t('mcpConsent.scopes.description') }}
									</p>
								</div>

								<ScopeGroupToggle
									v-for="scope in consent.requested_scopes"
									:key="scope.group"
									:scope="scope"
									:read="selection[scope.group]?.read === true"
									:edit="selection[scope.group]?.edit === true"
									:disabled="isActionDisabled"
									@update:read="setRead(scope.group, $event)"
									@update:edit="setEdit(scope.group, $event)"
								/>

								<p
									v-if="showSelectionError"
									class="text-destructive text-sm"
									role="alert"
									data-test="selection-error"
								>
									{{ t('mcpConsent.scopes.emptySelection') }}
								</p>
							</fieldset>

							<Alert v-if="hasEditSelected" variant="destructive" data-test="edit-warning">
								<Icon name="lucide:shield-alert" />
								<AlertTitle>{{ t('mcpConsent.writeWarning.title') }}</AlertTitle>
								<AlertDescription>{{ t('mcpConsent.writeWarning.description') }}</AlertDescription>
							</Alert>
						</CardContent>
						<CardFooter class="flex flex-col-reverse gap-3 sm:flex-row sm:justify-end">
							<Button
								:disabled="isActionDisabled"
								type="button"
								variant="outline"
								data-test="deny-button"
								@click="deny"
							>
								{{ t('mcpConsent.deny') }}
							</Button>
							<Button :disabled="isActionDisabled" type="submit" data-test="approve-button">
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
