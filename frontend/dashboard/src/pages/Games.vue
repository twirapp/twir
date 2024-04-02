<script setup lang="ts">

import { NGrid, NGridItem } from 'naive-ui';

import { responsiveCols } from '@/components/consants.js';
import EightBall from '@/components/games/8ball.vue';
import Duel from '@/components/games/duel.vue';
import RussianRoulette from '@/components/games/russianRoulette.vue';
</script>

<template>
	<div class="flex items-center justify-center max-w-[60vw] mx-auto my-0">
		<n-grid :cols="responsiveCols" :x-gap="16" :y-gap="16" responsive="screen">
			<n-grid-item :span="1">
				<EightBall />
			</n-grid-item>

			<n-grid-item :span="1">
				<RussianRoulette />
			</n-grid-item>

			<n-grid-item :span="1">
				<Duel />
			</n-grid-item>
		</n-grid>
	</div>
</template>
