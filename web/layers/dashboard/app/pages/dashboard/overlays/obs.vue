<script setup lang="ts">
import ObsOverlay from '~~/layers/dashboard/app/features/overlays/obs/page.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<ObsOverlay />
</template>
