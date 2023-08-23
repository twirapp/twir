<script setup lang="ts">
import { NAlert } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import List from '@/components/alerts/list.vue';

const { t } = useI18n();
</script>

<template>
	<n-alert type="info">
		<span v-html="t('alerts.info')" />
	</n-alert>
	<list />
</template>
