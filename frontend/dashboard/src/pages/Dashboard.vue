<script setup lang='ts'>
import { NGrid, NGridItem } from 'naive-ui';

import Bot from '@/components/dashboard/bot.vue';
import Stats from '@/components/dashboard/stats.vue';
</script>

<template>
	<Stats />
	<n-grid x-gap="12" :cols="3" style="padding: 24px;">
		<n-grid-item span="1">
			<Bot />
		</n-grid-item>
	</n-grid>
</template>

