<script setup lang="ts">
import VkCallback from '~~/layers/dashboard/features/integrations/pages/vk-callback.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<VkCallback />
</template>
