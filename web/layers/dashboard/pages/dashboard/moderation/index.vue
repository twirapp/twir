<script setup lang="ts">
import ModerationList from '~~/layers/dashboard/features/moderation/moderation.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<ModerationList />
</template>
