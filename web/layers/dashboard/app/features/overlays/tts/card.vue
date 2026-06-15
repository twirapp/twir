<script setup lang="ts">
import { useRouter } from 'vue-router'

import { useTTSOverlayApi } from '~~/layers/dashboard/app/api/overlays-tts'
import Card from '~~/layers/dashboard/app/components/overlays/card.vue'

const { t } = useI18n()
const api = useTTSOverlayApi()
const { data: settings, error, fetching } = api.useQueryTTS()

const router = useRouter()
</script>

<template>
	<card
		icon="twir-overlays:voice-message"
		title="TTS"
		:description="t('overlays.tts.description')"
		overlay-path="tts"
		:copy-disabled="!settings || !!error || !!fetching"
		@open-settings="router.push('/dashboard/overlays/tts')"
	>
	</card>
</template>

