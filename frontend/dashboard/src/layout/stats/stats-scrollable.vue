<script setup lang="ts">
import { NScrollbar, useThemeVars } from 'naive-ui'
import { computed } from 'vue'

import Stats from './stats.vue'

const themeVars = useThemeVars()
const blockColor = computed(() => themeVars.value.buttonColor2)
</script>

<template>
	<NScrollbar x-scrollable style="border-radius: 10px;">
		<div class="block scrollable-stats-block">
			<Stats />
		</div>
	</NScrollbar>
</template>

<style scoped>
.scrollable-stats-block {
	height: 45px;
	width: 100%;
	max-width: max-content;
	min-width: 40%;
	padding: unset;
	background-color: v-bind(blockColor);
	border-radius: 10px;
}

@media screen and (max-width: 768px) {
	.scrollable-stats-block {
		max-width: 100%;
		min-width: 100%;
	}
}
</style>
