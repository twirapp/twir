<script setup lang="ts">
import TtsOverlay from '~/layers/dashboard/features/overlays/tts/card.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<TtsOverlay />
</template>
