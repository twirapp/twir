<script setup lang="ts">
import { useProfile } from '~/api/user'
import Features from '~/components/landing/index/features/features.vue'
import Hero from '~/components/landing/index/hero/hero.vue'
import Integrations from '~/components/landing/index/integrations/integrations.vue'
import Stats from '~/components/landing/index/stats/stats.vue'
import Streamers from '~/components/landing/index/streamers/streamers.vue'
import Team from '~/components/landing/index/team/team.vue'

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
