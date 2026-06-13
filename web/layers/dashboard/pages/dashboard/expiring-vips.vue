<script setup lang="ts">
import ExpiringVips from '~~/layers/dashboard/features/expiring-vips/expiring-vips.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<ExpiringVips />
</template>
