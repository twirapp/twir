<script setup lang="ts">
import { useLocalStorage, useBreakpoints, breakpointsTailwind } from '@vueuse/core';
import {
	darkTheme,
	lightTheme,
	NLayout,
	NLayoutHeader,
	NLayoutContent,
	NLayoutSider,
	NConfigProvider,
	NMessageProvider,
	NNotificationProvider,
	NSpin,
} from 'naive-ui';
import { computed, ref, watch } from 'vue';
import { RouterView, useRouter } from 'vue-router';

import { useTheme } from '@/hooks/index.js';
import Header from '@/layout/header.vue';
import Sidebar from '@/layout/sidebar.vue';

const { theme } = useTheme();
const themeStyles = computed(() => theme.value === 'dark' ? darkTheme : lightTheme);

const isRouterReady = ref(false);
const router = useRouter();

router.isReady().finally(() => isRouterReady.value = true);

const breakPoints = useBreakpoints(breakpointsTailwind);
// If we are on a smaller than or equal to lg, we want the sidebar to collapse.
const smallerOrEqualLg = breakPoints.smallerOrEqual('lg');
// If we are on a smaller than or equal to md, we want the sidebar to hide and show hamburger menu with drawer.
const smallerOrEqualMd = breakPoints.smallerOrEqual('md');

const storedSidebarValue = useLocalStorage('twirSidebarIsCollapsed', false);

const toggleSidebar = () => {
	storedSidebarValue.value = !storedSidebarValue.value;
};

const isSidebarCollapsed = computed(() => {
	return storedSidebarValue.value;
});

watch(smallerOrEqualLg, (v) => {
	storedSidebarValue.value = v;
});
</script>

<template>
	<n-config-provider
		:theme="themeStyles"
		style="height: 100%"
		:breakpoints="{ xs: 0, s: 640, m: 1024, l: 1280, xl: 1536, xxl: 1920, '2xl': 2560 }"
	>
		<n-notification-provider :max="5">
			<n-message-provider>
				<n-layout style="height: 100%">
					<n-layout-header bordered style="height: var(--layout-header-height); width: 100%;">
						<Header :toggleSidebar="toggleSidebar" />
					</n-layout-header>

					<n-layout has-sider style="height: calc(100vh - var(--layout-header-height))">
						<n-layout-sider
							v-if="!smallerOrEqualMd"
							bordered
							collapse-mode="width"
							:collapsed-width="64"
							:width="240"
							show-trigger="arrow-circle"
							:native-scrollbar="false"
							:collapsed="isSidebarCollapsed"
							:show-collapsed-content="false"
							@update-collapsed="toggleSidebar"
						>
							<Sidebar :is-collapsed="isSidebarCollapsed" />
						</n-layout-sider>
						<n-layout-content>
							<div v-if="!isRouterReady" class="app-loader">
								<n-spin size="large" />
							</div>
							<router-view v-else v-slot="{ Component, route }">
								<transition :name="route.meta.transition as string || 'router'" mode="out-in">
									<div
										:key="route.path"
										:style="{
											padding: route.meta?.noPadding ? undefined: '24px'
										}"
									>
										<component :is="Component" />
									</div>
								</transition>
							</router-view>
						</n-layout-content>
					</n-layout>
				</n-layout>
			</n-message-provider>
		</n-notification-provider>
	</n-config-provider>
</template>

<style>
.router-enter-active,
.router-leave-active {
	transition: all 0.2s cubic-bezier(0, 0, 0.2, 1);
}

.router-enter-from,
.router-leave-to {
	opacity: 0;
	transform: scale(0.98);
}

.app-loader {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 100%;
}
</style>
