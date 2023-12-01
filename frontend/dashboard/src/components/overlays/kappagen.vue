<script setup lang="ts">
import { IconMoodWink } from '@tabler/icons-vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import { useKappaGenOverlayManager } from '@/api/index.js';
import Card from '@/components/overlays/card.vue';


const manager = useKappaGenOverlayManager();
const { data: settings, isError } = manager.getSettings();
const { t } = useI18n();

const router = useRouter();
</script>

<template>
	<card
		:icon="IconMoodWink"
		title="Kappagen"
		:description="t('overlays.kappagen.description')"
		overlay-path="kappagen"
		:copy-disabled="!settings || isError"
		@open-settings="router.push({ name: 'KappagenOverlay' })"
	>
	</card>
</template>
