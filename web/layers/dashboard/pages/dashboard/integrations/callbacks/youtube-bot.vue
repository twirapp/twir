<script setup lang="ts">
import YouTubeBotCallback from '~~/layers/dashboard/features/admin-panel/actions/ui/youtube-bot-callback.vue'

definePageMeta({ layout: 'popup', middleware: 'auth' })
</script>

<template>
	<YouTubeBotCallback />
</template>
