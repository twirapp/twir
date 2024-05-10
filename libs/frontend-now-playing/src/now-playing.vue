<script setup lang="ts">
import './assets/style.css'

import { useFontSource } from '@twir/fontsource'
import { ChannelOverlayNowPlayingPreset } from '@twir/types/api'
import { computed, watch } from 'vue'

import PresetAidenRedesign from './presets/aiden-redesign.vue'
import PresetSimpleLine from './presets/simple-line.vue'
import PresetTransparent from './presets/transparent.vue'

import type { Track } from './types.js'
import type { Settings } from '@twir/api/messages/overlays_now_playing/overlays_now_playing'

const props = defineProps<{
	settings: Settings
	track?: Track | null
}>()

const presetComponent = computed(() => {
	switch (props.settings.preset) {
		case ChannelOverlayNowPlayingPreset.TRANSPARENT:
			return PresetTransparent
		case ChannelOverlayNowPlayingPreset.AIDEN_REDESIGN:
			return PresetAidenRedesign
		case ChannelOverlayNowPlayingPreset.SIMPLE_LINE:
			return PresetSimpleLine
		default:
			return PresetTransparent
	}
})

const fontSource = useFontSource(false)

watch(() => [
	props.settings.fontFamily,
	props.settings.fontWeight,
], () => {
	fontSource.loadFont(
		props.settings.fontFamily,
		props.settings.fontWeight,
		'normal',
	)
}, { deep: true, immediate: true })

const fontFamily = computed(() => {
	return `"${props.settings.fontFamily}-${props.settings.fontWeight}-normal"`
})
</script>

<template>
	<component :is="presetComponent" :track="track" :settings="props.settings" />
</template>

<style>
.artist {
	font-family: v-bind(fontFamily);
	font-weight: v-bind('props.settings.fontWeight');
}

.name {
	font-family: v-bind(fontFamily);
	font-weight: v-bind('props.settings.fontWeight');
}

.image {
	flex-shrink: 0;
}
</style>
