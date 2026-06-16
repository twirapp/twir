<script setup lang="ts">

import { useRouter } from 'vue-router'

import Card from '~~/layers/dashboard/components/overlays/card.vue'

const { t } = useI18n()
const router = useRouter()
const localePath = useLocalePath()
</script>

<template>
	<Card
		icon="twir-overlays:broadcast"
		:icon-stroke="2"
		title="OBS"
		:description="t('overlays.obs.description')"
		overlay-path="obs"
		:showCopy="false"
		@open-settings="router.push(localePath('/dashboard/overlays/obs'))"
	>
	</Card>
</template>
