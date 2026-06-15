<script setup lang="ts">
import NightbotCallback from '~~/layers/dashboard/app/features/import/nightbot/nightbot-callback.vue'

definePageMeta({ layout: 'popup', middleware: 'auth' })
</script>

<template>
	<NightbotCallback />
</template>
