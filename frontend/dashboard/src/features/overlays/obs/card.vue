<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { useObsWebsocketApi } from '@/api/overlays-obs'
import BroadcastIcon from '@/assets/overlays/broadcast.svg?use'
import Card from '@/components/overlays/card.vue'

const { t } = useI18n()
const api = useObsWebsocketApi()
const { data: settings, error, fetching } = api.useQueryObsWebsocket()

const router = useRouter()
</script>

<template>
	<Card
		:icon="BroadcastIcon"
		:icon-stroke="2"
		title="OBS"
		:description="t('overlays.obs.description')"
		overlay-path="obs"
		:copy-disabled="!settings || !!error || !!fetching"
		@open-settings="router.push({ name: 'ObsOverlay' })"
	>
	</Card>
</template>
