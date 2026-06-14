<script setup lang="ts">
import BrbOverlay from '~~/layers/dashboard/features/overlays/brb/page.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<BrbOverlay />
</template>
