<script setup lang="ts">
import KappagenOverlay from '~/layers/dashboard/features/overlays/kappagen/card.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<KappagenOverlay />
</template>
