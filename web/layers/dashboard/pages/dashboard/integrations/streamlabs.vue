<script setup lang="ts">
import StreamlabsCallback from '~~/layers/dashboard/features/import/streamlabs/streamlabs-callback.vue'

definePageMeta({ layout: 'popup', middleware: 'auth' })
</script>

<template>
	<StreamlabsCallback />
</template>
