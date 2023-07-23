<script setup lang='ts'>
import { NGrid, NGridItem } from 'naive-ui';

import Bot from '@/components/dashboard/bot.vue';
import Chat from '@/components/dashboard/chat.vue';
import Stats from '@/components/dashboard/stats.vue';
</script>

<template>
	<Stats />
	<n-grid :x-gap="12" :y-gap="12" cols="1 sm:1 md:2 l:4" responsive="screen" style="padding: 24px; width: auto">
		<n-grid-item :span="2">
			<Chat />
		</n-grid-item>
		<n-grid-item :span="1">
			<Bot />
		</n-grid-item>
	</n-grid>
</template>

