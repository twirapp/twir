<script setup lang="ts">
import FaceitStatsOverlay from '~~/layers/dashboard/features/overlays/faceit-stats/builder.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<FaceitStatsOverlay />
</template>
