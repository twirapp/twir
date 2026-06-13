<script setup lang="ts">
import ValorantStatsOverlay from '~/layers/dashboard/features/overlays/valorant-stats/ui/card.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<ValorantStatsOverlay />
</template>
