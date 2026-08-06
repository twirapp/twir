<script setup lang="ts">
import { computed } from 'vue'

import type { Layer } from '../types'

interface Props {
	layer: Layer
}

const props = defineProps<Props>()

const textStyle = computed(() => ({
	color: props.layer.settings.textColor,
	fontFamily: props.layer.settings.textFontFamily,
	fontSize: `${props.layer.settings.textFontSize}px`,
	fontWeight: props.layer.settings.textFontWeight,
	textAlign: props.layer.settings.textAlign,
	whiteSpace: 'pre-wrap' as const,
	wordBreak: 'break-word' as const,
}))

const iframeStyle = computed(() => {
	// Mirror the runtime iframe-layer math: box stays at layer size, the iframe
	// itself is enlarged by 1/scale and then scaled down, so the preview matches
	// what OBS renders. pointer-events are disabled so the layer stays
	// draggable/selectable on the canvas (iframes swallow mouse events otherwise).
	const scale = props.layer.settings.iframeScale || 1
	return {
		width: `${100 / scale}%`,
		height: `${100 / scale}%`,
		transform: `scale(${scale})`,
		transformOrigin: 'top left',
	}
})

const youtubeEmbedUrl = computed(() => {
	const videoId = props.layer.settings.youtubeVideoId.trim()
	if (!videoId) return ''

	const params = new URLSearchParams({
		autoplay: props.layer.settings.youtubeAutoplay ? '1' : '0',
		mute: props.layer.settings.youtubeMuted ? '1' : '0',
		loop: props.layer.settings.youtubeLoop ? '1' : '0',
		controls: '0',
	})
	if (props.layer.settings.youtubeLoop) params.set('playlist', videoId)

	return `https://www.youtube-nocookie.com/embed/${encodeURIComponent(videoId)}?${params.toString()}`
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
			v-if="layer.settings.iframeUrl"
			:src="layer.settings.iframeUrl"
			:title="layer.name"
			:style="iframeStyle"
			class="pointer-events-none border-0"
		/>
		<div v-else class="flex h-full w-full items-center justify-center text-xs text-muted-foreground">
			Укажите URL виджета
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
			Укажите ID видео YouTube
		</div>
	</template>
	<template v-else-if="layer.type === 'EMOTE'">
		<img
			v-if="layer.settings.emoteUrl"
			:src="layer.settings.emoteUrl"
			alt="Emote"
			class="h-full w-full object-contain"
		/>
		<div v-else class="flex h-full w-full items-center justify-center text-xs text-muted-foreground">
			Выберите эмоцию
		</div>
	</template>
</template>
