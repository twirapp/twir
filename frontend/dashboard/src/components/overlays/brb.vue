<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import { useBeRightBackOverlayManager } from '@/api/index.js';
import KeyboardAndMouseIcon from '@/assets/overlays/keyboard-and-mouse.svg?use';
import Card from '@/components/overlays/card.vue';

const { t } = useI18n();
const manager = useBeRightBackOverlayManager();
const { data: settings, isError, isLoading } = manager.getSettings();

const router = useRouter();
</script>

<template>
	<card
		:icon="KeyboardAndMouseIcon"
		title="AFK overlay"
		:description="t('overlays.brb.description')"
		overlay-path="brb"
		:copy-disabled="!settings || isError || isLoading"
		@open-settings="router.push({ name: 'BrbOverlay' })"
	>
	</card>
</template>
