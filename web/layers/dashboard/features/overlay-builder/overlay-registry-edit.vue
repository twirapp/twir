<script setup lang="ts">
definePageMeta({ layout: 'dashboard', middleware: 'auth' })
const { t } = useI18n()
</script>

<template>
	<div>{{ t('overlayBuilder.overlayEdit') }}</div>
</template>
