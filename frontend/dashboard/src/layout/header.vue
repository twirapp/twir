<script setup lang="ts">
import { breakpointsTailwind, useBreakpoints, useEventListener } from '@vueuse/core';
import { useThemeVars } from 'naive-ui';
import { computed, ref, watch } from 'vue';
import { useRoute } from 'vue-router';

import ButtonDiscord from './buttons/buttonDiscord.vue';
import ButtonLogout from './buttons/buttonLogout.vue';
import ButtonPublicPage from './buttons/buttonPublicPage.vue';
import ButtonToggleTheme from './buttons/buttonToggleTheme.vue';
import DashboardsMenu from './dashboardsMenu.vue';
import Drawer from './drawer.vue';
import DropdownLanguage from './dropdowns/dropdownLanguage.vue';
import DropdownProfileOptions from './dropdowns/dropdownProfileOptions.vue';
import HamburgerMenu from './hamburgerMenu.vue';
import Logo from '../../public/TwirInCircle.svg?component';

import { useProfile } from '@/api';

defineProps<{
	toggleSidebar: () => void;
}>();

const route = useRoute();

const themeVars = useThemeVars();
const blockColor = computed(() => themeVars.value.buttonColor2);

const { data: profileData, isLoading: isProfileLoading } = useProfile();

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

	handleCloseDrawer();
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

			<div v-if="!smallerOrEqualMd" style="display: flex; gap: 12px;">
				<div class="block">
					<button-public-page />
					<button-discord />
				</div>

				<div class="block">
					<dropdown-language />
					<button-toggle-theme />
					<dropdown-profile-options :is-profile-loading="isProfileLoading" :avatar="profileData?.avatar" />
				</div>
			</div>

			<template v-if="smallerOrEqualMd">
				<hamburger-menu :is-open="isDrawerOpen" @click="handleToggleDrawer" />
				<drawer :show="isDrawerOpen">
					<div style="justify-content: space-between; padding: 8px;" class="block">
						<div style="display: flex; align-items: center; gap: 8px">
							<dropdown-language />
							<button-toggle-theme />
						</div>
						<button-logout :is-profile-loading="isProfileLoading" />
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
}

.block {
	background-color: v-bind(blockColor);
	display: flex;
	gap: 16px;
	padding: 16px;
	border-radius: 10px;
	align-items: center;
}
</style>
