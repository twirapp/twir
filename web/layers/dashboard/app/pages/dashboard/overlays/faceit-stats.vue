<script setup lang="ts">
import FaceitStatsOverlay from '~~/layers/dashboard/app/features/overlays/faceit-stats/builder.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<FaceitStatsOverlay />
</template>
