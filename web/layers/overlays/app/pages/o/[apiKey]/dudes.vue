<script setup lang="ts">
import DudesOverlay from '@twirapp/dudes-vue'
import { computed, onMounted, onUnmounted, watch } from 'vue'

import { assetsLoaderOptions, dudesSounds } from '~~/layers/overlays/app/composables/dudes/dudes-config'
import { useDudesIframe } from '~~/layers/overlays/app/composables/dudes/use-dudes-iframe'
import { useDudesSettings } from '~~/layers/overlays/app/composables/dudes/use-dudes-settings'
import { useDudesSocket } from '~~/layers/overlays/app/composables/dudes/use-dudes-socket'
import { useDudes } from '~~/layers/overlays/app/composables/dudes/use-dudes'
import { type ChatMessage, type ChatSettings, useChatTmi } from '~~/layers/overlays/app/composables/tmi/use-chat-tmi'
import { normalizeDisplayName } from '~~/layers/overlays/app/helpers'

definePageMeta({ layout: false })

const route = useRoute()
const { dudes, isDudeOverlayReady, createDude } = useDudes()
const { channelData, dudesSettings } = useDudesSettings()
const dudesSocketStore = useDudesSocket()
const iframe = useDudesIframe()

watch([isDudeOverlayReady, dudesSettings], ([isReady, settings]) => {
	if (!isReady || !settings || !dudes.value?.dudes) return
	dudes.value.dudes.updateSettings(settings.dudes)

	if (iframe.isIframe) {
		iframe.spawnIframeDude()
	}
})

async function onMessage(chatMessage: ChatMessage): Promise<void> {
	if (!dudes.value || chatMessage.type === 'system') return

	if (
		dudesSettings.value?.ignore.ignoreUsers &&
		dudesSettings.value.ignore.users.includes(chatMessage.senderId!)
	) {
		return
	}

	const name = normalizeDisplayName(chatMessage.senderDisplayName!, chatMessage.sender!)
	const dude = await createDude({
		userName: name,
		userId: chatMessage.senderId!,
		color: chatMessage.senderColor,
	})
	dude?.showMessage(chatMessage.chunks)
}

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: channelData.value?.channelId ?? '',
		channelName: channelData.value?.channelName ?? '',
		emotes: {
			isSmaller: true,
			bttv: true,
			ffz: true,
			sevenTv: true,
		},
		onMessage,
	}
})

const { destroy } = useChatTmi(chatSettings)

onMounted(async () => {
	const apiKey = route.params.apiKey as string
	const overlayId = route.query.id as string

	dudesSocketStore.connect(apiKey, overlayId)
	iframe.connect()
})

onUnmounted(() => {
	destroy()
	iframe.destroy()
})
</script>

<template>
	<DudesOverlay ref="dudes" :assets-loader-options="assetsLoaderOptions" :sounds="dudesSounds" />
</template>

<style>
body {
	overflow: hidden;
}
</style>