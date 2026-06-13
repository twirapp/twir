<script setup lang="ts">
import ObsOverlay from '~/layers/dashboard/features/overlays/obs/card.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<ObsOverlay />
</template>
