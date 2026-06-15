<script setup lang="ts">
;
import { useRouter } from 'vue-router';

import Card from '~~/layers/dashboard/components/overlays/card.vue';

const { t } = useI18n();

const router = useRouter();
const localePath = useLocalePath()
</script>

<template>
	<card
		icon="twir-overlays:dudes"
		title="Dudes"
		:description="t('overlays.dudes.description')"
		:show-copy="false"
		@open-settings="router.push(localePath('/dashboard/overlays/dudes'))"
	>
	</card>
</template>
