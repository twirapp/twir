<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import { useBeRightBackOverlayManager } from '@/api/index.js';
import IconBrb from '@/assets/icons/overlays/brb.svg?component';
import Card from '@/components/overlays/card.vue';

const { t } = useI18n();
const manager = useBeRightBackOverlayManager();
const { data: settings, isError, isLoading } = manager.getSettings();

const router = useRouter();
</script>

<template>
	<card
		:icon="IconBrb"
		title="AFK overlay"
		:description="t('overlays.brb.description')"
		overlay-path="brb"
		:copy-disabled="!settings || isError || isLoading"
		@open-settings="router.push({ name: 'BrbOverlay' })"
	>
	</card>
</template>
