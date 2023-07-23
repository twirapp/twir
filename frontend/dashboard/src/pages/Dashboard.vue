<script setup lang='ts'>
import { NGrid, NGridItem } from 'naive-ui';

import Bot from '@/components/dashboard/bot.vue';
</script>

<template>
	<n-grid x-gap="12" :cols="3">
		<n-grid-item span="1">
			<Bot />
		</n-grid-item>
	</n-grid>
</template>

