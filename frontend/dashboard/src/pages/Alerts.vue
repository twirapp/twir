<script setup lang="ts">
import { NAlert } from 'naive-ui';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useProfile } from '@/api/index.js';
import List from '@/components/alerts/list.vue';
import copyInput from '@/components/copyInput.vue';
import { storeToRefs } from 'pinia';

const { t } = useI18n();
const userProfile = storeToRefs(useProfile());
const overlayLink = computed(() => {
	return `${window.location.origin}/overlays/${userProfile.data?.value?.apiKey}/alerts`;
});
</script>

<template>
	<n-alert type="info" :title="t('alerts.info')">
		<span>{{ t('alerts.overlayLabel') }}</span>
		<copy-input type="password" :text="overlayLink" size="medium" />
	</n-alert>

	<list />
</template>
