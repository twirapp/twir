<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'

import { useTTSOverlayApi } from '@/api/overlays-tts'
import VoiceMessageIcon from '@/assets/overlays/voice-message.svg?use'
import Card from '@/components/overlays/card.vue'

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

