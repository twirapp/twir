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
} from 'naive-ui';
import { computed, ref, watch, onMounted } from 'vue';
import { RouterView } from 'vue-router';

import { useTheme } from '@/hooks/index.js';
import Header from '@/layout/header.vue';
import Sidebar from '@/layout/sidebar.vue';

const { theme } = useTheme();
const themeStyles = computed(() => theme.value === 'dark' ? darkTheme : lightTheme);


const breakPoints = useBreakpoints(breakpointsTailwind);
const smallerOrEqualLg = breakPoints.smallerOrEqual('lg');

const storedSidebarValue = useLocalStorage('twirSidebarIsCollapsed', false);
const toggleSidebar = () => storedSidebarValue.value = !storedSidebarValue.value;

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
            bordered
            collapse-mode="width"
            :collapsed-width="64"
            :width="240"
            :native-scrollbar="false"
            :collapsed="isSidebarCollapsed"
            :show-collapsed-content="false"
          >
            <Sidebar />
          </n-layout-sider>
          <n-layout-content content-style="padding: 24px; width: 100%">
            <router-view v-slot="{ Component, route }">
              <transition :name="route.meta.transition || 'fade'" mode="out-in">
                <div :key="route.name">
                  <component :is="Component" />
                </div>
              </transition>
            </router-view>
          </n-layout-content>
        </n-layout>
      </n-layout>
    </n-message-provider>
  </n-config-provider>
</template>


<style>
.fade-enter-active,
.fade-leave-active {
	transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
	opacity: 0;
}
</style>
