<script setup lang="ts">
import { NScrollbar } from 'naive-ui'

import Stats from './stats.vue'
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
}
</style>
