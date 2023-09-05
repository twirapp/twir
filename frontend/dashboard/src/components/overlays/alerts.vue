<script setup lang="ts">
import { IconBell } from '@tabler/icons-vue';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useProfile } from '@/api/index.js';
import Card from '@/components/overlays/card.vue';

const userProfile = useProfile();
const overlayLink = computed(() => {
	return `${window.location.origin}/overlays/${userProfile.data?.value?.apiKey}/alerts`;
});

const { t } = useI18n();
</script>

<template>
	<card
		:icon="IconBell"
		title="Alerts"
		:description="t('overlays.alerts.description')"
		:overlay-link="overlayLink"
		:show-settings="false"
	>
	</card>
</template>
