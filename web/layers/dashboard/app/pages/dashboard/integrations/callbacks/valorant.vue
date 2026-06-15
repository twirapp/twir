<script setup lang="ts">
import ValorantCallback from '~~/layers/dashboard/app/features/integrations/pages/valorant-callback.vue'

definePageMeta({ layout: 'popup', middleware: 'auth' })
</script>

<template>
	<ValorantCallback />
</template>
