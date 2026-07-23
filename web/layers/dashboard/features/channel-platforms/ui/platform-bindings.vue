<script setup lang="ts">
import type { Platform } from '~/gql/graphql.js'

import { toast } from 'vue-sonner'

import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Skeleton } from '@/components/ui/skeleton'

import { useChannelPlatforms } from '../composables/use-channel-platforms.js'
import PlatformBindingCard from './platform-binding-card.vue'

const {
	cards,
	fetching,
	error,
	connect: connectPlatform,
	disconnect: disconnectPlatform,
	setEnabled: setPlatformEnabled,
} = useChannelPlatforms()
const busyPlatform = ref<Platform | null>(null)

async function runAction(platform: Platform, action: () => Promise<unknown>) {
	busyPlatform.value = platform

	try {
		if (await action()) {
			toast.error('Unable to update platform binding')
		}
	} catch {
		toast.error('Unable to update platform binding')
	} finally {
		busyPlatform.value = null
	}
}

function connect(platform: Platform) {
	return runAction(platform, () => connectPlatform(platform))
}

function disconnect(platform: Platform) {
	return runAction(platform, () => disconnectPlatform(platform))
}

function setEnabled(platform: Platform, enabled: boolean) {
	return runAction(platform, () => setPlatformEnabled(platform, enabled))
}
</script>

<template>
	<section class="flex flex-col gap-4">
		<div class="flex flex-col gap-1">
			<h2 class="text-xl font-semibold tracking-tight">Platform bindings</h2>
			<p class="text-sm text-muted-foreground">
				Connect the accounts your bot can use across streaming platforms.
			</p>
		</div>

		<Alert
			v-if="error"
			variant="destructive"
		>
			<Icon name="lucide:circle-alert" />
			<AlertTitle>Unable to load platform bindings</AlertTitle>
			<AlertDescription>{{ error.message }}</AlertDescription>
		</Alert>

		<div
			v-if="fetching && cards.length === 0"
			class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3"
		>
			<Skeleton
				v-for="index in 3"
				:key="index"
				class="h-64"
			/>
		</div>

		<div
			v-else
			class="grid gap-4 sm:grid-cols-2 xl:grid-cols-3"
		>
			<PlatformBindingCard
				v-for="card in cards"
				:key="card.platform"
				:platform="card.platform"
				:presentation="card.presentation"
				:capabilities="card.capabilities"
				:binding="card.binding"
				:busy="busyPlatform === card.platform"
				@connect="connect"
				@disconnect="disconnect"
				@set-enabled="setEnabled"
			/>
		</div>
	</section>
</template>
