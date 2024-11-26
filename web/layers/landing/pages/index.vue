<script setup lang="ts">
import { useProfile } from '~/layers/landing/api/user'
import Features from '~/layers/landing/components/index/features/features.vue'
import Hero from '~/layers/landing/components/index/hero/hero.vue'
import Integrations from '~/layers/landing/components/index/integrations/integrations.vue'
import Stats from '~/layers/landing/components/index/stats/stats.vue'
import Streamers from '~/layers/landing/components/index/streamers/streamers.vue'
import Team from '~/layers/landing/components/index/team/team.vue'

await useProfile()

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
		<Streamers />
		<Team />
	</div>
</template>
