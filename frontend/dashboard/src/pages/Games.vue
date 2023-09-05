<script setup lang="ts">

import { NGrid, NGridItem } from 'naive-ui';

import EightBall from '@/components/games/8ball.vue';
</script>

<template>
	<div
		style="
			display: flex;
			align-items: center;
			justify-content: center;
			max-width: 1000px;
			margin: 0 auto;
		"
	>
		<n-grid cols="1 m:2" :x-gap="16" :y-gap="16" responsive="screen">
			<n-grid-item :span="1">
				<EightBall />
			</n-grid-item>
		</n-grid>
	</div>
</template>

<style scoped>

</style>
