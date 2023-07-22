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
	NDrawer,
	NDrawerContent,
} from 'naive-ui';
import { computed, watch } from 'vue';
import { RouterView } from 'vue-router';

import { useTheme } from '@/hooks/index.js';
import Header from '@/layout/header.vue';
import Sidebar from '@/layout/sidebar.vue';

const { theme } = useTheme();
const themeStyles = computed(() => theme.value === 'dark' ? darkTheme : lightTheme);

const breakPoints = useBreakpoints(breakpointsTailwind);
const smallerOrEqualLg = breakPoints.smallerOrEqual('lg');

const storedSidebarValue = useLocalStorage('twirSidebarIsCollapsed', false);
const storedDrawerValue = useLocalStorage('twirDrawerIsCollapsed', false);

const toggleSidebar = () => {
	if (smallerOrEqualLg.value) {
		storedDrawerValue.value = !storedDrawerValue.value;
	} else {
		storedSidebarValue.value = !storedSidebarValue.value;
	}
};

const isSidebarCollapsed = computed(() => {
	return storedSidebarValue.value;
});

watch(smallerOrEqualLg, (v) => {
	storedSidebarValue.value = v;
});

</script>

<template>
	<n-config-provider :theme="themeStyles" style="height: 100%">
		<n-message-provider>
			<n-layout style="height: 100%">
				<n-layout-header bordered style="height: 43px;">
					<Header :toggleSidebar="toggleSidebar" />
				</n-layout-header>
				<n-layout has-sider style="height: calc(100vh - 43px)">
					<n-layout-sider
						v-if="!smallerOrEqualLg"
						bordered
						collapse-mode="width"
						:collapsed-width="64"
						:width="240"
						:native-scrollbar="false"
						:collapsed="isSidebarCollapsed"
						:show-collapsed-content="false"
					>
						<Sidebar :is-collapsed="isSidebarCollapsed" />
					</n-layout-sider>
					<n-drawer v-else v-model:show="storedDrawerValue" placement="left">
						<n-drawer-content body-content-style="padding: 0px">
							<Sidebar :is-collapsed="isSidebarCollapsed" />
						</n-drawer-content>
					</n-drawer>
					<n-layout-content content-style="padding-left: 24px; padding-right:24px; width: 100%">
						<router-view v-slot="{ Component, route }">
							<!-- TODO: THIS TRANSITION TRIGGERING WHEN WE OPENING DRAWER(MOBILES) -->
							<!-- <transition :name="route.meta.transition as string || 'router'" mode="out-in"> -->
							<div :key="route.name ?? Date.now()">
								<component :is="Component" />
							</div>
							<!-- </transition> -->
						</router-view>
					</n-layout-content>
				</n-layout>
			</n-layout>
		</n-message-provider>
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
</style>
