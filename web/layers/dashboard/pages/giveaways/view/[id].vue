<script setup lang="ts">
import GiveawayView from '~/layers/dashboard/features/giveaways/giveaway-view.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })
</script>

<template>
	<GiveawayView />
</template>
