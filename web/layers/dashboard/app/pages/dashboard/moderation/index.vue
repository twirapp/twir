<script setup lang="ts">
import ModerationList from '~~/layers/dashboard/app/features/moderation/moderation.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<ModerationList />
</template>
