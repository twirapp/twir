<script setup lang="ts">
import type { ChannelPlatformBinding, Platform } from '~/gql/graphql.js'

import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
	AlertDialogTrigger,
} from '@/components/ui/alert-dialog'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'

const props = withDefaults(
	defineProps<{
		platform: Platform
		presentation: { label: string; icon: string }
		capabilities: { name: string }[]
		binding: ChannelPlatformBinding | null
		busy?: boolean
	}>(),
	{
		busy: false,
	},
)

const emit = defineEmits<{
	connect: [platform: Platform]
	disconnect: [platform: Platform]
	setEnabled: [platform: Platform, enabled: boolean]
}>()

const description = computed(() => {
	if (props.binding) return `Connected as ${props.binding.platformDisplayName}`

	return `Connect your ${props.presentation.label} account to enable it in Twir.`
})

function connect() {
	emit('connect', props.platform)
}

function disconnect() {
	emit('disconnect', props.platform)
}

function setEnabled(enabled: boolean) {
	emit('setEnabled', props.platform, enabled)
}
</script>

<template>
	<Card class="flex h-full flex-col">
		<CardHeader class="flex flex-row items-start gap-3">
			<Icon
				:name="presentation.icon"
				class="size-5 shrink-0 text-muted-foreground"
			/>
			<div class="flex min-w-0 flex-1 flex-col gap-1">
				<CardTitle>{{ presentation.label }}</CardTitle>
				<CardDescription>{{ description }}</CardDescription>
			</div>
		</CardHeader>

		<CardContent class="flex flex-1 flex-col gap-4">
			<div
				v-if="binding"
				class="flex items-center gap-3"
			>
				<Avatar class="size-10">
					<AvatarImage
						v-if="binding.platformAvatar"
						:src="binding.platformAvatar"
					/>
					<AvatarFallback>{{ binding.platformDisplayName.charAt(0) }}</AvatarFallback>
				</Avatar>
				<div class="flex min-w-0 flex-1 flex-col gap-1">
					<p class="truncate font-medium">{{ binding.platformDisplayName }}</p>
					<p class="truncate text-sm text-muted-foreground">@{{ binding.platformLogin }}</p>
				</div>
			</div>
			<p
				v-else
				class="text-sm text-muted-foreground"
			>
				No account connected.
			</p>

			<div
				v-if="capabilities.length"
				class="flex flex-wrap gap-2"
			>
				<Badge
					v-for="capability in capabilities"
					:key="capability.name"
					variant="secondary"
				>
					{{ capability.name }}
				</Badge>
			</div>

			<div
				v-if="binding"
				class="flex items-center justify-between gap-4"
			>
				<label
					:for="`platform-enabled-${platform}`"
					class="text-sm font-medium"
				>
					Enable bot
				</label>
				<Switch
					:id="`platform-enabled-${platform}`"
					:model-value="binding.enabled"
					:disabled="busy"
					@update:model-value="setEnabled"
				/>
			</div>
		</CardContent>

		<CardFooter class="mt-auto">
			<Button
				v-if="!binding"
				class="w-full"
				:disabled="busy"
				@click="connect"
			>
				Connect
			</Button>
			<AlertDialog v-else>
				<AlertDialogTrigger as-child>
					<Button
						class="w-full"
						variant="destructive"
						:disabled="busy"
					>
						Disconnect
					</Button>
				</AlertDialogTrigger>
				<AlertDialogContent>
					<AlertDialogHeader>
						<AlertDialogTitle>Disconnect {{ presentation.label }}?</AlertDialogTitle>
						<AlertDialogDescription>
							The bot will stop using this connected account.
						</AlertDialogDescription>
					</AlertDialogHeader>
					<AlertDialogFooter>
						<AlertDialogCancel>Cancel</AlertDialogCancel>
						<AlertDialogAction
							:disabled="busy"
							@click="disconnect"
						>
							Disconnect
						</AlertDialogAction>
					</AlertDialogFooter>
				</AlertDialogContent>
			</AlertDialog>
		</CardFooter>
	</Card>
</template>
