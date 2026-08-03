<script setup lang="ts">
import { Accordion } from '@/components/ui/accordion'
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
import { useMcpConsentScopeSelection } from './mcp-consent-selection.js'
import ScopeGroupToggle from './scope-group-toggle.vue'

type ConsentScreenState = 'loading' | 'ready' | 'expired' | 'permission' | 'network'

const { t } = useI18n()
const route = useRoute()
const mcpConsentApi = createMcpConsentApi($fetch)
const attempt = computed(() => parseMcpConsentAttempt(route.query.attempt))
const consent = ref<McpConsent | null>(null)
const screenState = ref<ConsentScreenState>('loading')
const isSubmitting = ref(false)
const hasCompleted = ref(false)
const {
	selection,
	showSelectionError,
	hasEditSelected,
	initializeSelection,
	setRead,
	setEdit,
	hasEveryRequestedActionSelected,
	toggleAll,
} = useMcpConsentScopeSelection()

const selectedDashboardName = computed(() => consent.value?.channel_id ?? '')

const isActionDisabled = computed(() => isSubmitting.value || hasCompleted.value)

const allRequestedActionsSelected = computed(() => consent.value !== null && hasEveryRequestedActionSelected(consent.value.requested_scopes))

function setScreenState(result: Awaited<ReturnType<typeof mcpConsentApi.getMcpConsent>>): void {
	switch (result.kind) {
		case 'success':
			consent.value = result.data
			initializeSelection(result.data.requested_scopes)
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

function toggleAllRequestedActions(): void {
	if (consent.value === null || isActionDisabled.value) return

	toggleAll(consent.value.requested_scopes)
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
	<div class="bg-background flex min-h-svh w-full">
		<div class="m-auto flex w-full max-w-lg flex-col gap-6 px-4 py-8">
			<div class="flex flex-col gap-1 text-center">
				<h1 class="text-xl font-semibold">{{ t('mcpConsent.title') }}</h1>
				<p class="text-muted-foreground text-sm">{{ t('mcpConsent.description') }}</p>
			</div>

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
					<CardContent class="flex flex-col gap-4">
						<dl class="flex flex-col divide-y rounded-lg border text-sm">
							<div class="flex items-center justify-between gap-3 px-3 py-2">
								<dt class="text-muted-foreground shrink-0">{{ t('mcpConsent.client') }}</dt>
								<dd class="truncate font-medium">{{ consent.client.name }}</dd>
							</div>
							<div class="flex items-center justify-between gap-3 px-3 py-2">
								<dt class="text-muted-foreground shrink-0">{{ t('mcpConsent.dashboard') }}</dt>
								<dd class="truncate font-medium">{{ selectedDashboardName }}</dd>
							</div>
							<div
								v-if="consent.client.uri"
								class="flex items-center justify-between gap-3 px-3 py-2"
							>
								<dt class="text-muted-foreground shrink-0">{{ t('mcpConsent.clientUri') }}</dt>
								<dd class="truncate font-medium">{{ consent.client.uri }}</dd>
							</div>
						</dl>

						<fieldset class="flex flex-col gap-3">
							<legend class="sr-only">{{ t('mcpConsent.scopes.title') }}</legend>
							<div class="flex flex-wrap items-start justify-between gap-3">
								<div class="flex flex-col gap-1">
									<h3 class="text-sm font-medium">{{ t('mcpConsent.scopes.title') }}</h3>
									<p class="text-muted-foreground text-sm">
										{{ t('mcpConsent.scopes.description') }}
									</p>
								</div>
								<Button
									:disabled="isActionDisabled"
									type="button"
									variant="outline"
									size="sm"
									data-test="bulk-selection-button"
									@click="toggleAllRequestedActions"
								>
									{{
										allRequestedActionsSelected
											? t('mcpConsent.scopes.deselectAll')
											: t('mcpConsent.scopes.selectAll')
									}}
								</Button>
							</div>

							<Accordion type="multiple" class="w-full rounded-lg border px-3">
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
							</Accordion>

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
	</div>
</template>
