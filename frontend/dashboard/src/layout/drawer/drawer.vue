<script setup lang="ts">
import { NDivider, NDrawer, NDrawerContent } from 'naive-ui'

import DrawerHeader from './drawer-header.vue'
import ButtonImportCommands from '../buttons/button-import-commands.vue'
import MenuDashboards from '../menu-dashboards.vue'
import NavigationResources from '../navigation/navigation-resources.vue'
import Navigation from '../navigation/navigation.vue'

import FeedbackButton from '@/layout/feedback/feedback-button.vue'

defineProps<{
	show: boolean
}>()
</script>

<template>
	<NDrawer
		:show="show"
		:show-mask="false"
		close-on-esc
		placement="right"
		width="100%"
		style="top: var(--layout-header-height)"
		height="calc(100vh - var(--layout-header-height))"
	>
		<NDrawerContent header-style="padding: 8px">
			<DrawerHeader />

			<MenuDashboards is-drawer />

			<NDivider class="!m-0 !mt-2.5" />

			<NavigationResources />

			<NDivider class="!m-0 !mt-2 !mb-0.5" />

			<Navigation />

			<ButtonImportCommands is-drawer />
			<FeedbackButton is-drawer />
		</NDrawerContent>
	</NDrawer>
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
