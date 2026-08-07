<script setup lang="ts">
import { Platform } from '~/gql/graphql.js'

import { useStreamBackground } from '../composables/useStreamBackground'
import { useTwitchEmbedPlayer } from '../composables/useTwitchEmbedPlayer'

// Rendered OUTSIDE the zoom-transformed canvas subtree: the Twitch player
// refuses to play when any ancestor has a CSS transform, so the background
// replicates the scaled canvas rect in unscaled coordinates.
interface Props {
	width: number
	height: number
	offsetX: number
	offsetY: number
}

const props = defineProps<Props>()

const { t } = useI18n()
const { preference, selectedBinding, streamPreviewSrc, showUnsupportedPreview } = useStreamBackground()

const twitchContainer = ref<HTMLElement>()
const twitchLogin = computed(() => {
	if (!preference.value.enabled || selectedBinding.value?.platform !== Platform.Twitch) return null
	return selectedBinding.value.platformLogin
})
const autoResume = computed(() => preference.value.enabled)

const { gateLifted } = useTwitchEmbedPlayer(twitchContainer, twitchLogin, autoResume)

const showPreview = computed(() => twitchLogin.value !== null || streamPreviewSrc.value !== null || showUnsupportedPreview.value)

const backgroundStyle = computed(() => ({
	left: '50%',
	top: '50%',
	width: `${props.width}px`,
	height: `${props.height}px`,
	marginLeft: `${-props.width / 2 + props.offsetX}px`,
	marginTop: `${-props.height / 2 + props.offsetY}px`,
	zIndex: gateLifted.value ? 50 : 0,
}))
</script>

<template>
	<div v-if="showPreview" class="pointer-events-none absolute left-1/2 top-1/2 overflow-hidden" :style="backgroundStyle">
		<div v-if="twitchLogin" ref="twitchContainer" class="size-full" />
		<iframe v-else-if="streamPreviewSrc" class="pointer-events-none size-full border-0" :src="streamPreviewSrc" :title="t('overlayBuilder.streamBackground.previewTitle')" allow="autoplay" />
		<div v-else class="flex size-full items-center justify-center bg-background/80 p-4 text-center text-xs text-muted-foreground">{{ t('overlayBuilder.streamBackground.previewNotSupported') }}</div>
	</div>
</template>
