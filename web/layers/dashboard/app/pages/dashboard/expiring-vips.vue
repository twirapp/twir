<script setup lang="ts">
import ExpiringVips from '~~/layers/dashboard/app/features/expiring-vips/expiring-vips.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<ExpiringVips />
</template>
