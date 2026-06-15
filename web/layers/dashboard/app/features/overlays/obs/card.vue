<script setup lang="ts">
import { useRouter } from 'vue-router'

import Card from '~~/layers/dashboard/app/components/overlays/card.vue'

const { t } = useI18n()
const router = useRouter()
</script>

<template>
	<Card
		icon="twir-overlays:broadcast"
		:icon-stroke="2"
		title="OBS"
		:description="t('overlays.obs.description')"
		overlay-path="obs"
		:showCopy="false"
		@open-settings="router.push('/dashboard/overlays/obs')"
	>
	</Card>
</template>
