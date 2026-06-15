<script setup lang="ts">

import { useRouter } from 'vue-router'

import { useTTSOverlayApi } from '~~/layers/dashboard/api/overlays-tts'
import Card from '~~/layers/dashboard/components/overlays/card.vue'

const { t } = useI18n()
const api = useTTSOverlayApi()
const { data: settings, error, fetching } = api.useQueryTTS()

const router = useRouter()
const localePath = useLocalePath()
</script>

<template>
	<card
		icon="twir-overlays:voice-message"
		title="TTS"
		:description="t('overlays.tts.description')"
		overlay-path="tts"
		:copy-disabled="!settings || !!error || !!fetching"
		@open-settings="router.push(localePath('/dashboard/overlays/tts'))"
	>
	</card>
</template>

