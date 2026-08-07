<script setup lang="ts">
import { useStreamBackground } from '../composables/useStreamBackground'

const { t } = useI18n()
const { streamPreviewSrc, showUnsupportedPreview } = useStreamBackground()
</script>

<template>
	<div v-if="streamPreviewSrc || showUnsupportedPreview" class="pointer-events-none absolute inset-0 overflow-hidden" :style="{ zIndex: 0 }">
		<iframe v-if="streamPreviewSrc" class="pointer-events-none size-full border-0" :src="streamPreviewSrc" :title="t('overlayBuilder.streamBackground.previewTitle')" allow="autoplay" />
		<div v-else class="flex size-full items-center justify-center bg-background/80 p-4 text-center text-xs text-muted-foreground">{{ t('overlayBuilder.streamBackground.previewNotSupported') }}</div>
	</div>
</template>
