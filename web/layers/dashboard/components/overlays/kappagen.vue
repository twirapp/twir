<script setup lang="ts">
import { IconMoodWink } from '@tabler/icons-vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import Card from '@/components/overlays/card.vue';

const { t } = useI18n();

const router = useRouter();
</script>

<template>
	<card
		:icon="IconMoodWink"
		:icon-stroke="1"
		title="Kappagen"
		:description="t('overlays.kappagen.description')"
		overlay-path="kappagen"
		:show-copy="false"
		@open-settings="router.push({ name: 'KappagenOverlay' })"
	>
	</card>
</template>
