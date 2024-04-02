<script setup lang="ts">
import { NDivider, NScrollbar } from 'naive-ui';

import ImportCommandsButton from './buttons/import-commands-button.vue';
import DashboardsMenu from './dashboardsMenu.vue';
import Navigation from './navigation.vue';

import FeedbackButton from '@/layout/feedback/feedback-button.vue';
</script>

<template>
	<div class="flex flex-col justify-between h-[calc(100vh_-_var(--layout-header-height))] w-full">
		<dashboards-menu />

		<n-divider class="!m-0" />

		<n-scrollbar trigger="none">
			<navigation />
		</n-scrollbar>

		<import-commands-button />
		<feedback-button />
	</div>
</template>
