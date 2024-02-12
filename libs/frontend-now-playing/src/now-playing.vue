<script setup lang="ts">
import type { Settings } from '@twir/api/messages/overlays_now_playing/overlays_now_playing';
import { ChannelOverlayNowPlayingPreset } from '@twir/types/api';
import { computed } from 'vue';

import PresetAidenRedesign from './presets/aiden-redesign.vue';
import PresetTransparent from './presets/transparent.vue';
import type { Track } from './types.js';

const props = defineProps<{
	settings: Settings
	track?: Track | null
}>();

const presetComponent = computed(() => {
	switch (props.settings.preset) {
		case ChannelOverlayNowPlayingPreset.TRANSPARENT:
			return PresetTransparent;
		case ChannelOverlayNowPlayingPreset.AIDEN_REDESIGN:
			return PresetAidenRedesign;
		default:
			return PresetTransparent;
	}
});
</script>

<template>
	<component :is="presetComponent" :track="track" />
</template>

