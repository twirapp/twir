<script setup lang="ts">
import { NDrawer, NDrawerContent } from 'naive-ui';

import Navigation from './navigation.vue';

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
		<n-drawer-content body-content-style="padding: 12px" header-style="padding: 8px">
			<template #header>
				<slot />
			</template>
			<navigation :is-collapsed="false" />
		</n-drawer-content>
	</n-drawer>
</template>

<style scoped>
:deep(.n-drawer-header__main) {
	width: 100%;
}
</style>
