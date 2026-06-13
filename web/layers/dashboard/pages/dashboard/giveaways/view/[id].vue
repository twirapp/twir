<script setup lang="ts">
import GiveawayView from '~~/layers/dashboard/features/giveaways/giveaways.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<GiveawayView />
</template>
