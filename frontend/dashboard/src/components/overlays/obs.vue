<script setup lang="ts">
import { NModal } from 'naive-ui';
import { ref, computed } from 'vue';

import Settings from './obs/settings.vue';

import { useProfile, useObsOverlayManager } from '@/api/index.js';
import BroadcastIcon from '@/assets/icons/overlays/obs.svg?component';
import Card from '@/components/overlays/card.vue';

const isModalOpened = ref(false);
const obsManager = useObsOverlayManager();
const { data: obsSettings } = obsManager.getSettings();

const userProfile = useProfile();
const overlayLink = computed(() => {
	if (!obsSettings.value?.serverAddress) return;

	return `${window.location.origin}/overlays/${userProfile.data?.value?.apiKey}/obs`;
});
</script>

<template>
	<card
		:icon="BroadcastIcon"
		title="OBS"
		description="This overlay used for connect TwirApp with your obs. It gives opportunity to bot manage your sources, scenes, audio sources on events."
		:overlay-link="overlayLink"
		@open-settings="isModalOpened = true"
	>
	</card>

	<n-modal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="OBS"
		content-style="padding: 10px; width: 100%"
		style="width: 500px; max-width: calc(100vw - 40px)"
	>
		<Settings />
	</n-modal>
</template>
