<script setup lang="ts">
import DiscordIntegration from '~~/layers/dashboard/features/integrations/pages/discord-settings.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<DiscordIntegration />
</template>
