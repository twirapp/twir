<script setup lang="ts">
import ValorantCallback from '~~/layers/dashboard/features/integrations/ui/valorant/callback.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<ValorantCallback />
</template>
