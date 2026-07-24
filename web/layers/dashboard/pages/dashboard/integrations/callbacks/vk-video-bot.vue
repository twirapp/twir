<script setup lang="ts">
import VKVideoBotCallback from '~~/layers/dashboard/features/admin-panel/actions/ui/vk-video-bot-callback.vue'

definePageMeta({ layout: 'popup', middleware: 'auth' })
</script>

<template>
	<VKVideoBotCallback />
</template>
