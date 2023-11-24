<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import Settings from './brb/settings.vue';

import { useBeRightBackOverlayManager } from '@/api/index.js';
import IconBrb from '@/assets/icons/overlays/brb.svg?component';
import Card from '@/components/overlays/card.vue';

const { t } = useI18n();
const isModalOpened = ref(false);
const manager = useBeRightBackOverlayManager();
const { data: settings, isError, isLoading } = manager.getSettings();
</script>

<template>
	<card
		:icon="IconBrb"
		title="Be right back (afk)"
		:description="t('overlays.brb.description')"
		overlay-path="brb"
		:copy-disabled="!settings || isError || isLoading"
		@open-settings="isModalOpened = true"
	>
	</card>

	<Settings :show-settings="isModalOpened" @close="isModalOpened = false" />
</template>
