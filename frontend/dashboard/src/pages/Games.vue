<script setup lang="ts">
import { NGrid, NGridItem } from 'naive-ui'

import { responsiveCols } from '@/components/consants.js'
import EightBall from '@/components/games/8ball.vue'
import Duel from '@/components/games/duel.vue'
import RussianRoulette from '@/components/games/russianRoulette.vue'
import Seppuku from '@/components/games/seppuku.vue'
</script>

<template>
	<div class="flex items-center justify-center max-w-[60vw] mx-auto my-0">
		<NGrid :cols="responsiveCols" :x-gap="16" :y-gap="16" responsive="screen">
			<NGridItem :span="1">
				<EightBall />
			</NGridItem>

			<NGridItem :span="1">
				<RussianRoulette />
			</NGridItem>

			<NGridItem :span="1">
				<Duel />
			</NGridItem>

			<NGridItem :span="1">
				<Seppuku />
			</NGridItem>
		</NGrid>
	</div>
</template>
