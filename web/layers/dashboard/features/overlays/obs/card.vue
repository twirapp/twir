<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import BroadcastIcon from '~~/layers/dashboard/assets/overlays/broadcast.svg'
import Card from '~~/layers/dashboard/components/overlays/card.vue'

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
