<script setup lang="ts">
import DiscordCallback from '~/layers/dashboard/features/integrations/ui/discord/callback.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<DiscordCallback />
</template>
