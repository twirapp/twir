<script setup lang="ts">
import KappagenOverlay from '~~/layers/dashboard/features/overlays/kappagen/kappagen.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<KappagenOverlay />
</template>
