<script setup lang="ts">
import { NModal } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import Settings from './obs/settings.vue';

import { useObsOverlayManager } from '@/api/index.js';
import BroadcastIcon from '@/assets/icons/overlays/obs.svg?component';
import Card from '@/components/overlays/card.vue';

const isModalOpened = ref(false);
const obsManager = useObsOverlayManager();
const { data: obsSettings, isError } = obsManager.getSettings();

const { t } = useI18n();
</script>

<template>
	<card
		:icon="BroadcastIcon"
		title="OBS"
		:description="t('overlays.obs.description')"
		overlay-path="obs"
		:copy-disabled="!obsSettings || isError"
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
