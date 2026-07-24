<script setup lang="ts">
import Platforms from '~~/layers/dashboard/features/platforms/platforms.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<Platforms />
</template>
