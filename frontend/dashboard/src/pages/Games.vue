<script setup lang="ts">
import EightBall from '@/components/games/8ball.vue'
import RussianRoulette from '@/components/games/russianRoulette.vue'
import Seppuku from '@/components/games/seppuku.vue'
import Duel from '@/features/games/duel.vue'
import Voteban from '@/features/games/voteban.vue'
import PageLayout from '@/layout/page-layout.vue'
</script>

<template>
	<PageLayout>
		<template #title>
			Games
		</template>

		<template #content>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<EightBall />
				<RussianRoulette />
				<Duel />
				<Seppuku />
				<Voteban />
			</div>
		</template>
	</PageLayout>
</template>
