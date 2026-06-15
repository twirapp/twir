<script setup lang="ts">
import BrbOverlay from '~~/layers/dashboard/app/features/overlays/brb/page.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<BrbOverlay />
</template>
