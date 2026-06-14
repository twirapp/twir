<script setup lang="ts">
import ValorantStatsOverlay from '~~/layers/dashboard/features/overlays/valorant-stats/builder.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<ValorantStatsOverlay />
</template>
