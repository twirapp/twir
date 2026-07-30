<script setup lang="ts">
import { toast } from 'vue-sonner'

import {
	useMutationYouTubeBotSetupLink,
	useMutationYouTubeBotSetupStatus,
	youtubeBotSetupBroadcastChannelName,
} from '~~/layers/dashboard/api/admin/actions'
import ActionConfirm from '@/components/ui/action-confirm/ActionConfirm.vue'
import Button from '@/components/ui/button/Button.vue'
import Card from '@/components/ui/card/Card.vue'
import CardContent from '@/components/ui/card/CardContent.vue'
import CardDescription from '@/components/ui/card/CardDescription.vue'
import CardHeader from '@/components/ui/card/CardHeader.vue'
import CardTitle from '@/components/ui/card/CardTitle.vue'
import { Platform } from '~/gql/graphql.js'
import { PLATFORM_META } from '~/utils/platforms.js'

const youtubeMeta = PLATFORM_META[Platform.Youtube]

const setupLinkMutation = useMutationYouTubeBotSetupLink()
const setupStatusMutation = useMutationYouTubeBotSetupStatus()

const configured = ref(false)
const isReplacing = ref(false)
const isStarting = ref(false)
const isStatusLoading = ref(true)
const refreshChannel = shallowRef<BroadcastChannel | null>(null)

async function refreshSetupStatus() {
	isStatusLoading.value = true
	const result = await setupStatusMutation.executeMutation({})
	isStatusLoading.value = false

	if (result.error) {
		toast.error(result.error.message)
		return
	}

	configured.value = result.data?.youtubeBotSetupStatus ?? false
}

async function startSetup() {
	isStarting.value = true
	const result = await setupLinkMutation.executeMutation({})
	isStarting.value = false

	if (result.error) {
		toast.error(result.error.message)
		return
	}

	const setupLink = result.data?.youtubeBotSetupLink
	if (!setupLink) {
		toast.error('Unable to start YouTube bot setup')
		return
	}

	window.open(setupLink, 'youtube-bot-setup', 'popup')
}

function requestSetup() {
	if (configured.value) {
		isReplacing.value = true
		return
	}

	void startSetup()
}

onMounted(() => {
	const channel = new BroadcastChannel(youtubeBotSetupBroadcastChannelName)
	channel.onmessage = (event) => {
		if (event.data === 'refresh') void refreshSetupStatus()
	}
	refreshChannel.value = channel
	void refreshSetupStatus()
})

onBeforeUnmount(() => {
	refreshChannel.value?.close()
})
</script>

<template>
	<Card>
		<CardHeader>
			<div class="flex items-center gap-2">
				<Icon :name="youtubeMeta.icon" class="size-5" :class="youtubeMeta.colorClass" />
				<CardTitle>YouTube bot</CardTitle>
			</div>
			<CardDescription>
				Configure the global bot account used by every YouTube channel binding.
			</CardDescription>
		</CardHeader>
		<CardContent class="flex items-center justify-between gap-4">
			<p class="text-sm text-muted-foreground">
				{{ isStatusLoading ? 'Checking setup status...' : configured ? 'Connected' : 'Not connected' }}
			</p>
			<Button :disabled="isStatusLoading || isStarting" @click="requestSetup">
				{{ isStarting ? 'Opening...' : configured ? 'Replace' : 'Connect' }}
			</Button>
		</CardContent>
	</Card>

	<ActionConfirm
		v-model:open="isReplacing"
		confirm-text="Replace the configured YouTube bot for every YouTube channel binding?"
		@confirm="startSetup"
	/>
</template>
