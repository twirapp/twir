<script setup lang="ts">
import TtsOverlay from '~~/layers/dashboard/app/features/overlays/tts/page.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<TtsOverlay />
</template>
