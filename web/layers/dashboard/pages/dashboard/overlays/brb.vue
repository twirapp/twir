<script setup lang="ts">
import BrbOverlay from '~~/layers/dashboard/features/overlays/brb/page.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<BrbOverlay />
</template>
