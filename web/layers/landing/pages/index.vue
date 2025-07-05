<script setup lang="ts">
import Features from '~~/layers/landing/components/index/features/features.vue'
import Hero from '~~/layers/landing/components/index/hero/hero.vue'
import Integrations from '~~/layers/landing/components/index/integrations/integrations.vue'
import Stats from '~~/layers/landing/components/index/stats/stats.vue'

definePageMeta({
	layout: 'landing',
})
</script>

<template>
	<div class="text-white">
		<Hero />
		<Stats />
		<Features />
		<Integrations />
		<!--		<Team /> -->
	</div>
</template>
