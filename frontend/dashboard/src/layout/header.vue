<script setup lang="ts">
import { breakpointsTailwind, useBreakpoints, useEventListener } from '@vueuse/core';
import { useThemeVars } from 'naive-ui';
import { computed, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import ButtonDiscord from './buttons/buttonDiscord.vue';
import ButtonPublicPage from './buttons/buttonPublicPage.vue';
import ButtonToggleTheme from './buttons/buttonToggleTheme.vue';
import DashboardsMenu from './dashboardsMenu.vue';
import Drawer from './drawer.vue';
import DropdownLanguage from './dropdowns/dropdownLanguage.vue';
import DropdownProfile from './dropdowns/dropdownProfile.vue';
import HamburgerMenu from './hamburgerMenu.vue';
import Stats from './stats.vue';
import Logo from '../../public/TwirInCircle.svg?component';


defineProps<{
	toggleSidebar: () => void;
}>();

const route = useRoute();

const themeVars = useThemeVars();
const blockColor = computed(() => themeVars.value.buttonColor2);

const breakPoints = useBreakpoints(breakpointsTailwind);
const smallerOrEqualMd = breakPoints.smallerOrEqual('md');

const isDrawerOpen = ref(false);
const handleCloseDrawer = () => isDrawerOpen.value = false;
const handleToggleDrawer = () => isDrawerOpen.value = !isDrawerOpen.value;

watch(smallerOrEqualMd, (v) => {
	if (!v && isDrawerOpen.value) {
		handleCloseDrawer();
	}
});

watch(() => route.fullPath, () => {
	if (!isDrawerOpen.value) return;

	handleCloseDrawer();
});

useEventListener('keydown', (ev) => {
	if (ev.code !== 'Escape' || !isDrawerOpen.value) return;

	// TODO: Don`t close if dropdown-profile-options is open
	// handleCloseDrawer();
});
</script>

<template>
	<div class="header">
		<div class="content">
			<div>
				<a href="/" class="logo">
					<Logo style="width: 36px; height: 36px; display: flex" />
				</a>

				<DashboardsMenu />
			</div>

			<div v-if="!smallerOrEqualMd" style="flex-grow: 1">
				<div class="block" style="width: 100%; max-width: max-content; min-width: 40%;">
					<Stats />
				</div>
			</div>

			<div v-if="!smallerOrEqualMd" style="display: flex; gap: 12px;">
				<div class="block">
					<button-public-page />
					<button-discord />
				</div>

				<div class="block">
					<dropdown-language />
					<button-toggle-theme />
					<dropdown-profile />
				</div>
			</div>

			<template v-if="smallerOrEqualMd">
				<hamburger-menu :is-open="isDrawerOpen" @click="handleToggleDrawer" />
				<drawer :show="isDrawerOpen">
					<div class="drawerSlot">
						<div class="drawerSlotBlocks">
							<div class="drawerBlock">
								<button-public-page />
								<button-discord />
							</div>

							<div class="drawerBlock">
								<dropdown-language />
								<button-toggle-theme />
							</div>
						</div>

						<dropdown-profile />
					</div>
				</drawer>
			</template>
		</div>
	</div>
</template>

<style scoped>
.header {
	height: 100%;
	padding-top: 10px;
	padding-left: 10px;
	padding-right: 10px;
}

.content .logo {
	display: flex;
	align-items: center;
	text-decoration: none;
	padding-right: 14px;
}

.content > div {
	display: flex;
	justify-content: flex-start;
	gap: 5px;
	height: 45px;
}

.content {
	display: flex;
	justify-content: space-between;
	align-items: center;
	height: 45px;
	column-gap: 12px;
}

.block {
	background-color: v-bind(blockColor);
	display: flex;
	gap: 16px;
	padding: 16px;
	border-radius: 10px;
	align-items: center;
}

.drawerBlock {
	background-color: v-bind(blockColor);
	display: flex;
	align-items: center;
	gap: 8px;
	padding: 6px 16px;
	border-radius: 10px;
}

.drawerSlot {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 8px;
}

.drawerSlotBlocks {
	display: flex;
	align-items: center;
	gap: 6px;
	margin-right: 6px;
}

@media screen and (max-width: 520px) {
	.drawerSlot {
		align-items: flex-start;
	}

	.drawerSlotBlocks {
		flex-direction: column;
	}
}
</style>
