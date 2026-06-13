<script setup lang="ts">
import DiscordCallback from '~~/layers/dashboard/features/integrations/pages/discord-callback.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<DiscordCallback />
</template>
