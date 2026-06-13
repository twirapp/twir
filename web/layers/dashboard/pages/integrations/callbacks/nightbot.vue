<script setup lang="ts">
import NightbotCallback from '~/layers/dashboard/features/integrations/ui/nightbot/callback.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<NightbotCallback />
</template>
