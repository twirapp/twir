<script setup lang="ts">
import DiscordIntegration from '~~/layers/dashboard/features/integrations/ui/discord/settings.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<DiscordIntegration />
</template>
