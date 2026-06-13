<script setup lang="ts">
import ModerationEdit from '~/layers/dashboard/features/moderation/moderation-edit.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<ModerationEdit />
</template>
