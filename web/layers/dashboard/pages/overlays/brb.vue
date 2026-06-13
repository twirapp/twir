<script setup lang="ts">
import BrbOverlay from '~/layers/dashboard/features/overlays/brb/card.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<BrbOverlay />
</template>
