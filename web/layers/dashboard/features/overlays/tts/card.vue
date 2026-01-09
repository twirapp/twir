<script setup lang="ts">

import { useRouter } from 'vue-router'

import { useTTSOverlayApi } from '#layers/dashboard/api/overlays-tts'
import VoiceMessageIcon from '~/assets/overlays/voice-message.svg?use'
import Card from '#layers/dashboard/components/overlays/card.vue'

const { t } = useI18n()
const api = useTTSOverlayApi()
const { data: settings, error, fetching } = api.useQueryTTS()

const router = useRouter()
</script>

<template>
	<card
		:icon="VoiceMessageIcon"
		title="TTS"
		:description="t('overlays.tts.description')"
		overlay-path="tts"
		:copy-disabled="!settings || !!error || !!fetching"
		@open-settings="router.push({ name: 'TTSOverlay' })"
	>
	</card>
</template>
