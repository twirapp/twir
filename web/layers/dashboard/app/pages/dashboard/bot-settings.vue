<script setup lang="ts">
import BotSettings from '~~/layers/dashboard/app/features/bot-settings/bot-settings.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<BotSettings />
</template>
