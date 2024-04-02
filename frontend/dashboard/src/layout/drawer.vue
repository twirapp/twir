<script setup lang="ts">
import { NDivider, NDrawer, NDrawerContent } from 'naive-ui';

import ImportCommandsButton from './buttons/import-commands-button.vue';
import DashboardsMenu from './dashboardsMenu.vue';
import Navigation from './navigation.vue';

import FeedbackButton from '@/layout/feedback/feedback-button.vue';

defineProps<{
	show: boolean
}>();
</script>

<template>
	<n-drawer
		:show="show"
		:show-mask="false"
		close-on-esc
		placement="left"
		width="100%"
		style="top: var(--layout-header-height)"
		height="calc(100vh - var(--layout-header-height))"
	>
		<n-drawer-content header-style="padding: 8px">
			<template #header>
				<slot />
			</template>

			<dashboards-menu is-drawer />

			<n-divider class="!m-0 !mt-2.5" />

			<navigation />

			<import-commands-button is-drawer />
			<feedback-button is-drawer />
		</n-drawer-content>
	</n-drawer>
</template>

<style scoped>
:deep(.n-drawer-body-content-wrapper) {
	padding: 12px !important;
}

:deep(.n-drawer-header__main) {
	width: 100%;
}

@media screen and (max-width: 520px) {
	:deep(.n-drawer-body-content-wrapper) {
		padding: 4px 0 !important;
	}
}
</style>
