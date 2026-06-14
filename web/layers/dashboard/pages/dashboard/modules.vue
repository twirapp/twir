<script setup lang="ts">
import Modules from '~~/layers/dashboard/features/modules/modules.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<Modules />
</template>
