<script setup lang="ts">
import GiveawaysPagesView from '~~/layers/dashboard/features/giveaways/pages/giveaways-pages-view.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<GiveawaysPagesView />
</template>
