<script setup lang="ts">
import StreamStatsOverlay from '~~/layers/dashboard/features/overlays/stream-stats/page.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<StreamStatsOverlay />
</template>
