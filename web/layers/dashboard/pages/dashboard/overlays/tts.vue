<script setup lang="ts">
import TtsOverlay from '~~/layers/dashboard/features/overlays/tts/page.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<TtsOverlay />
</template>
