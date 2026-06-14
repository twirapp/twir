<script setup lang="ts">

import Card from '~~/layers/dashboard/components/overlays/card.vue';

const { t } = useI18n();
</script>

<template>
	<card
		:icon="'tabler:icon-bell'"
		title="Alerts"
		:description="t('overlays.alerts.description')"
		overlay-path="alerts"
		:show-settings="false"
	>
	</card>
</template>
