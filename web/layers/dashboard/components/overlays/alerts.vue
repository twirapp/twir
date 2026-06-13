<script setup lang="ts">
import { IconBell } from '@tabler/icons-vue';
import { useI18n } from 'vue-i18n';

import Card from '@/components/overlays/card.vue';

const { t } = useI18n();
</script>

<template>
	<card
		:icon="IconBell"
		title="Alerts"
		:description="t('overlays.alerts.description')"
		overlay-path="alerts"
		:show-settings="false"
	>
	</card>
</template>
