<script setup lang="ts">
import { NGrid, NGridItem } from 'naive-ui';

import { responsiveCols } from '@/components/consants.js';
import Alerts from '@/components/overlays/alerts.vue';
import Chat from '@/components/overlays/chat.vue';
import Kappagen from '@/components/overlays/kappagen.vue';
import OBS from '@/components/overlays/obs.vue';
import TTS from '@/components/overlays/tts.vue';
</script>

<template>
	<div
		style="
			display: flex;
			align-items: center;
			justify-content: center;
			max-width: 60vw;
			margin: 0 auto;
		"
	>
		<n-grid :cols="responsiveCols" :x-gap="16" :y-gap="16" responsive="screen">
			<n-grid-item :span="1">
				<TTS />
			</n-grid-item>
			<n-grid-item :span="1">
				<OBS />
			</n-grid-item>
			<n-grid-item :span="1">
				<Alerts />
			</n-grid-item>
			<n-grid-item :span="1">
				<Chat />
			</n-grid-item>
			<n-grid-item :span="1">
				<Kappagen />
			</n-grid-item>
		</n-grid>
	</div>
</template>
