<script setup lang="ts">
import { IconMoodWink } from '@tabler/icons-vue';
import { NModal } from 'naive-ui';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import Settings from './kappagen/settings.vue';

import { useChatOverlayManager } from '@/api/index.js';
import Card from '@/components/overlays/card.vue';

const isModalOpened = ref(false);
const chatManager = useChatOverlayManager();
const { data: settings, isError } = chatManager.getSettings();
const { t } = useI18n();
</script>

<template>
	<card
		:icon="IconMoodWink"
		title="Kappagen"
		:description="t('overlays.kappagen.description')"
		overlay-path="kappagen"
		:copy-disabled="!settings || isError"
		@open-settings="isModalOpened = true"
	>
	</card>

	<n-modal
		v-model:show="isModalOpened"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		title="Kappagen"
		content-style="padding: 10px; width: 100%"
		style="width: 70vw; max-width: calc(100vw - 40px)"
	>
		<Settings />
	</n-modal>
</template>
