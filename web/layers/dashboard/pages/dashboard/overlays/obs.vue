<script setup lang="ts">
import ObsOverlay from '~~/layers/dashboard/features/overlays/obs/page.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<ObsOverlay />
</template>
