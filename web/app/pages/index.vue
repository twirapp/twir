<script setup lang="ts">
import { profileQuery, userInvalidateQueryKey } from '~/api/user'
import Features from '~/pages/index/features/features.vue'
import Hero from '~/pages/index/hero/hero.vue'
import Integrations from '~/pages/index/integrations/integrations.vue'
import Stats from '~/pages/index/stats/stats.vue'
import Team from '~/pages/index/team/team.vue'

definePageMeta({
	layout: 'landing',
})

await useQuery({
	query: profileQuery.query,
	variables: {},
	context: {
		key: profileQuery.key,
		additionalTypenames: [userInvalidateQueryKey],
	},
})
</script>

<template>
	<div class="text-white">
		<Hero />
		<Stats />
		<Features />
		<Integrations />
		<Team />
	</div>
</template>
