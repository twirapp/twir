<script setup lang="ts">
import { useRouter } from 'vue-router'
import Card from '~~/layers/dashboard/components/overlays/card.vue'

import { useStreamStatsOverlayApi } from './api'

const { t } = useI18n()
const api = useStreamStatsOverlayApi()
const { data: settings, error, fetching } = api.useQueryStreamStats()

const router = useRouter()
const localePath = useLocalePath()
</script>

<template>
	<Card
		icon="lucide:chart-column"
		:title="t('overlays.streamStats.title')"
		:description="t('overlays.streamStats.description')"
		overlay-path="stream-stats"
		:copy-disabled="!settings || !!error || !!fetching"
		@open-settings="router.push(localePath('/dashboard/overlays/stream-stats'))"
	>
	</Card>
</template>
