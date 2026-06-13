<script setup lang="ts">
import GiveawaysList from '~~/layers/dashboard/features/giveaways/giveaways.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<GiveawaysList />
</template>
