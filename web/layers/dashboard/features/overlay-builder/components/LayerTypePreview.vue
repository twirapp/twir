<script setup lang="ts">
import { toRef } from 'vue'

import type { Layer } from '../types'
import { useLayerTypePreview } from '../composables/useLayerTypePreview'

interface Props {
	layer: Layer
}

const props = defineProps<Props>()
const { t } = useI18n()

const { textStyle, iframeStyle, youtubeEmbedUrl } = useLayerTypePreview(toRef(props, 'layer'))

const iframePreviewUrl = computed(() => {
	const raw = props.layer.settings.iframeUrl
	if (!raw) return ''

	try {
		const url = new URL(raw)
		if (url.origin !== window.location.origin) return raw
		url.searchParams.set('preview', '1')
		return url.toString()
	} catch {
		return raw
	}
})
</script>

<template>
	<div v-if="layer.type === 'TEXT'" class="flex h-full w-full items-center justify-center overflow-hidden px-2" :style="textStyle">
		{{ layer.settings.textContent }}
	</div>
	<video
		v-else-if="layer.type === 'VIDEO'"
		:src="layer.settings.videoUrl || undefined"
		:muted="layer.settings.videoMuted"
		:loop="layer.settings.videoLoop"
		autoplay
		playsinline
		class="pointer-events-none h-full w-full object-cover"
	/>
	<template v-else-if="layer.type === 'IFRAME'">
		<iframe
			v-if="iframePreviewUrl"
			:src="iframePreviewUrl"
			:title="layer.name"
			:style="iframeStyle"
			class="pointer-events-none border-0"
		/>
		<div v-else class="flex h-full w-full items-center justify-center text-xs text-muted-foreground">
			{{ t('overlayBuilder.preview.widgetUrl') }}
		</div>
	</template>
	<template v-else-if="layer.type === 'YOUTUBE'">
		<iframe
			v-if="youtubeEmbedUrl"
			:src="youtubeEmbedUrl"
			:title="layer.name"
			class="pointer-events-none h-full w-full border-0"
		/>
		<div v-else class="flex h-full w-full items-center justify-center text-xs text-muted-foreground">
			{{ t('overlayBuilder.preview.youtubeId') }}
		</div>
	</template>
	<template v-else-if="layer.type === 'EMOTE'">
		<img
			v-if="layer.settings.emoteUrl"
			:src="layer.settings.emoteUrl"
			:alt="t('overlayBuilder.preview.emoteAlt')"
			class="h-full w-full object-contain"
		/>
		<div v-else class="flex h-full w-full items-center justify-center text-xs text-muted-foreground">
			{{ t('overlayBuilder.preview.selectEmote') }}
		</div>
	</template>
</template>
