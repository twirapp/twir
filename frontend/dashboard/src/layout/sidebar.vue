<script setup lang="ts">
import { NScrollbar } from 'naive-ui';

import Navigation from './navigation.vue';

import FeedbackButton from '@/layout/feedback/feedback-button.vue';


defineProps<{
	isCollapsed: boolean
}>();

</script>

<template>
	<div class="sidebar">
		<n-scrollbar trigger="none">
			<navigation :is-collapsed="isCollapsed" />
		</n-scrollbar>

		<feedback-button :is-collapsed="isCollapsed" />
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
