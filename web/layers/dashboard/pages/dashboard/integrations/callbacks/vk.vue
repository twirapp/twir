<script setup lang="ts">
import VkCallback from '~~/layers/dashboard/features/integrations/pages/vk-callback.vue'

definePageMeta({ layout: 'popup', middleware: 'auth' })
</script>

<template>
	<VkCallback />
</template>
