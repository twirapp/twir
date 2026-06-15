<script setup lang="ts">
import FilesList from '~~/layers/dashboard/app/components/files/files.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<FilesList mode="list" />
</template>
