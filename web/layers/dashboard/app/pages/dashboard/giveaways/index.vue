<script setup lang="ts">
import GiveawaysList from '~~/layers/dashboard/app/features/giveaways/giveaways.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<GiveawaysList />
</template>
