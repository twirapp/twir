<script setup lang="ts">
import { toRef } from 'vue'

import type { Layer } from '../types'
import { useLayerTypePreview } from '../composables/useLayerTypePreview'

interface Props {
	layer: Layer
}

const props = defineProps<Props>()

const { textStyle, iframeStyle, youtubeEmbedUrl } = useLayerTypePreview(toRef(props, 'layer'))
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
