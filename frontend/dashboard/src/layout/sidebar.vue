<script setup lang="ts">
import { NDivider, NScrollbar } from 'naive-ui'

import ButtonImportCommands from './buttons/button-import-commands.vue'
import FeedbackButton from './feedback/feedback-button.vue'
import MenuDashboards from './menu-dashboards.vue'
import Navigation from './navigation/navigation.vue'
</script>

<template>
	<div class="flex flex-col justify-between h-[calc(100vh_-_var(--layout-header-height))] w-full">
		<MenuDashboards />

		<NDivider class="!m-0" />

		<NScrollbar trigger="none">
			<Navigation />
		</NScrollbar>

		<ButtonImportCommands />
		<FeedbackButton />
	</div>
</template>
