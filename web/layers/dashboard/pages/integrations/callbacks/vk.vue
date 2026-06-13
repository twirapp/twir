<script setup lang="ts">
import VkCallback from '~/layers/dashboard/features/integrations/ui/vk/callback.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<VkCallback />
</template>
