<script setup lang="ts">
import { Settings } from '@twir/api/messages/overlays_now_playing/overlays_now_playing';

import type { Track } from '../types.js';

defineProps<{
	track?: Track | null
	settings: Settings
}>();
</script>

<template>
	<div v-if="track" class="spotify">
		<div class="info">
			<span class="name">{{ track.title }}</span>
			<span>–</span>
			<div class="artist">
				{{ track.artist }}
			</div>
		</div>
	</div>
</template>

<style scoped>
body {
	display: flex;
}

.spotify {
	display: flex;
	align-items: center;
	column-gap: 16px;
	background-color: v-bind('settings.backgroundColor');

	/* Font */
	font-family: Inter;
	font-size: 24px;
	color: #fff;
	padding-top: 4px;
	padding-bottom: 4px;
	padding-left: 8px;
	padding-right: 8px;
	border-radius: 8px;
}

.info {
	display: flex;
	flex-direction: row;
	gap: 4px;
}
</style>
