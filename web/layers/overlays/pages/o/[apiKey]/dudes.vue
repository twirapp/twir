<script setup lang="ts">
import DudesOverlay from '@twirapp/dudes-vue'
import { computed, onMounted, onUnmounted, watch } from 'vue'
import { useRoute } from 'vue-router'

import { assetsLoaderOptions, dudesSounds } from '~/layers/overlays/composables/dudes/dudes-config.ts'
import { useDudesIframe } from '~/layers/overlays/composables/dudes/use-dudes-iframe.ts'
import { useDudesSettings } from '~/layers/overlays/composables/dudes/use-dudes-settings.ts'
import { useDudesSocket } from '~/layers/overlays/composables/dudes/use-dudes-socket.ts'
import { useDudes } from '~/layers/overlays/composables/dudes/use-dudes.ts'
import { type ChatMessage, type ChatSettings, useChatTmi } from '~/layers/overlays/composables/tmi/use-chat-tmi.ts'
import { normalizeDisplayName } from '~/layers/overlays/helpers.ts'

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