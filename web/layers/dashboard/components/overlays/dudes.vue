<script setup lang="ts">
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import DudesIcon from '@/assets/overlays/dudes.svg?use';
import Card from '@/components/overlays/card.vue';

const { t } = useI18n();

const router = useRouter();
</script>

<template>
	<card
		:icon="DudesIcon"
		title="Dudes"
		:description="t('overlays.dudes.description')"
		:show-copy="false"
		@open-settings="router.push({ name: 'DudesOverlay' })"
	>
	</card>
</template>
