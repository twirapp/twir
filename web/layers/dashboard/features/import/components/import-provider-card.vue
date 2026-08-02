<script setup lang="ts">
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar/index.js'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert/index.js'
import { Badge } from '@/components/ui/badge/index.js'
import { Button } from '@/components/ui/button/index.js'
import {
	Card,
	CardContent,
	CardDescription,
	CardFooter,
	CardHeader,
	CardTitle,
} from '@/components/ui/card/index.js'
import { Dialog, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog/index.js'
import DialogOrSheet from '~~/layers/dashboard/components/dialog-or-sheet.vue'

interface ProviderAccount {
	userName?: string | null
	avatar?: string | null
}

const props = withDefaults(
	defineProps<{
		title: string
		icon: string
		description: string
		account?: ProviderAccount | null
		authLink?: string | null
		isLoading?: boolean
		canManageIntegration: boolean
		importAvailable?: boolean
		donationDescription?: string
		unavailableDescription?: string
	}>(),
	{
		account: null,
		authLink: null,
		isLoading: false,
		importAvailable: true,
		donationDescription: '',
		unavailableDescription: '',
	}
)

defineEmits<{
	logout: []
}>()

defineSlots<{
	settings?: () => unknown
}>()

const { t } = useI18n()
const settingsOpen = ref(false)
const connected = computed(() => Boolean(props.account?.userName))

function authenticate() {
	if (!props.authLink) return
	window.open(props.authLink, `Twir connect ${props.title}`, 'width=800,height=600')
}
</script>

<template>
	<Card class="flex h-full flex-col">
		<CardHeader>
			<div class="flex items-start justify-between gap-4">
				<div class="flex items-center gap-3">
					<div class="flex size-11 shrink-0 items-center justify-center rounded-lg bg-muted">
						<Icon :name="icon" class="size-7" />
					</div>
					<div class="flex min-w-0 flex-col gap-1">
						<CardTitle>{{ title }}</CardTitle>
						<Badge :variant="connected ? 'secondary' : 'outline'">
							{{ connected ? t('imports.status.connected') : t('imports.status.disconnected') }}
						</Badge>
					</div>
				</div>
				<Badge v-if="!importAvailable" variant="outline">
					{{ t('imports.status.unavailable') }}
				</Badge>
			</div>
			<CardDescription>{{ description }}</CardDescription>
		</CardHeader>

		<CardContent class="flex grow flex-col gap-4">
			<div v-if="connected" class="flex items-center gap-3 rounded-lg border p-3">
				<Avatar>
					<AvatarImage v-if="account?.avatar" :src="account.avatar" :alt="account.userName ?? title" />
					<AvatarFallback>{{ account?.userName?.slice(0, 2).toUpperCase() }}</AvatarFallback>
				</Avatar>
				<div class="min-w-0">
					<p class="text-xs text-muted-foreground">{{ t('imports.status.connectedAs') }}</p>
					<p class="truncate font-medium">{{ account?.userName }}</p>
				</div>
			</div>

			<Alert v-if="donationDescription">
				<Icon name="lucide:heart-handshake" />
				<AlertTitle>{{ t('imports.status.donationsEnabled') }}</AlertTitle>
				<AlertDescription>{{ donationDescription }}</AlertDescription>
			</Alert>

			<Alert v-if="!importAvailable">
				<Icon name="lucide:info" />
				<AlertTitle>{{ t('imports.status.importUnavailable') }}</AlertTitle>
				<AlertDescription>{{ unavailableDescription }}</AlertDescription>
			</Alert>
		</CardContent>

		<CardFooter class="mt-auto flex flex-wrap gap-2">
			<Button
				v-if="importAvailable && connected && $slots.settings"
				variant="secondary"
				:disabled="isLoading"
				data-test="provider-settings"
				@click="settingsOpen = true"
			>
				<Icon name="lucide:settings" data-icon="inline-start" />
				{{ t('imports.actions.settings') }}
			</Button>
			<Button
				:variant="connected ? 'destructive' : 'default'"
				:disabled="!canManageIntegration || isLoading || (!connected && !authLink)"
				data-test="provider-auth"
				@click="connected ? $emit('logout') : authenticate()"
			>
				<Icon v-if="isLoading" name="lucide:loader-circle" data-icon="inline-start" class="animate-spin" />
				<Icon v-else :name="connected ? 'lucide:log-out' : 'lucide:log-in'" data-icon="inline-start" />
				{{ connected ? t('imports.actions.disconnect') : t('imports.actions.connect') }}
			</Button>
		</CardFooter>
	</Card>

	<Dialog v-if="importAvailable && $slots.settings" v-model:open="settingsOpen">
		<DialogOrSheet class="max-w-4xl">
			<DialogHeader>
				<DialogTitle>{{ t('imports.settings.title', { provider: title }) }}</DialogTitle>
				<DialogDescription>{{ t('imports.settings.description') }}</DialogDescription>
			</DialogHeader>
			<slot name="settings" />
		</DialogOrSheet>
	</Dialog>
</template>
