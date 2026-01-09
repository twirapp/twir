<script setup lang="ts">

import { useRouter } from 'vue-router'

import BroadcastIcon from '~/assets/overlays/broadcast.svg?use'
import Card from '#layers/dashboard/components/overlays/card.vue'

const { t } = useI18n()
const router = useRouter()
</script>

<template>
	<Card
		:icon="BroadcastIcon"
		:icon-stroke="2"
		title="OBS"
		:description="t('overlays.obs.description')"
		overlay-path="obs"
		:showCopy="false"
		@open-settings="router.push({ name: 'ObsOverlay' })"
	>
	</Card>
</template>
