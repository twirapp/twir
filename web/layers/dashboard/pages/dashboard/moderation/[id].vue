<script setup lang="ts">
import ModerationEdit from '~~/layers/dashboard/features/moderation/moderation-form.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<ModerationEdit />
</template>
