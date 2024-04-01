<script setup lang="ts">
import { NScrollbar } from 'naive-ui';

import ImportCommandsButton from './buttons/import-commands-button.vue';
import Navigation from './navigation.vue';

import FeedbackButton from '@/layout/feedback/feedback-button.vue';
</script>

<template>
	<div class="sidebar">
		<n-scrollbar trigger="none">
			<navigation />
		</n-scrollbar>

		<import-commands-button />
		<feedback-button />
	</div>
</template>

<style scoped>
.sidebar {
	display: flex;
	flex-direction: column;
	justify-content: space-between;
	height: calc(100vh - var(--layout-header-height));
	width: 100%;
}
</style>
