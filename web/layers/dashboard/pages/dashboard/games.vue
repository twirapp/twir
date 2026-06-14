<script setup lang="ts">
import EightBall from '~~/layers/dashboard/components/games/8ball.vue'
import RussianRoulette from '~~/layers/dashboard/components/games/russianRoulette.vue'
import Seppuku from '~~/layers/dashboard/components/games/seppuku.vue'
import Duel from '~~/layers/dashboard/features/games/duel.vue'
import Voteban from '~~/layers/dashboard/features/games/voteban.vue'
import PageLayout from '~~/layers/dashboard/layout/page-layout.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth', noPadding: true })
</script>

<template>
	<PageLayout>
		<template #title> Games </template>

		<template #content>
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2">
				<EightBall />
				<RussianRoulette />
				<Duel />
				<Seppuku />
				<Voteban />
			</div>
		</template>
	</PageLayout>
</template>
