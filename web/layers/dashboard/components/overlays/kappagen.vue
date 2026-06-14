<script setup lang="ts">
import { useRouter } from 'vue-router';

import Card from '~~/layers/dashboard/components/overlays/card.vue';

const { t } = useI18n();

const router = useRouter();
</script>

<template>
	<card
		:icon="'tabler:icon-mood-wink'"
		:icon-stroke="1"
		title="Kappagen"
		:description="t('overlays.kappagen.description')"
		overlay-path="kappagen"
		:show-copy="false"
		@open-settings="router.push('/dashboard/overlays/kappagen')"
	>
	</card>
</template>
