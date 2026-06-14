<script setup lang="ts">
import DiscordCallback from '~~/layers/dashboard/features/integrations/pages/discord-callback.vue'

definePageMeta({ layout: 'popup', middleware: 'auth' })
</script>

<template>
	<DiscordCallback />
</template>
