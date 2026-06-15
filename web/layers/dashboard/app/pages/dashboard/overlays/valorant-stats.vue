<script setup lang="ts">
import ValorantStatsOverlay from '~~/layers/dashboard/app/features/overlays/valorant-stats/builder.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<ValorantStatsOverlay />
</template>
