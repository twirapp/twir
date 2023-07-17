<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core';
import {
	darkTheme,
	lightTheme,
	NLayout,
	NLayoutHeader,
	NLayoutContent,
	NLayoutSider,
	NConfigProvider,
} from 'naive-ui';
import { computed, ref } from 'vue';
import { RouterView } from 'vue-router';

import { useTheme } from '@/hooks/index.js';
import Header from '@/layout/header.vue';
import Sidebar from '@/layout/sidebar.vue';

const { theme } = useTheme();
const themeStyles = computed(() => theme.value === 'dark' ? darkTheme : lightTheme);

const sidebarCollapsed = useLocalStorage('twirIsSidebarCollapsed', false);
const toggleSidebar = () => sidebarCollapsed.value = !sidebarCollapsed.value;
</script>

<template>
  <n-config-provider :theme="themeStyles" style="height: 100%">
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
          :collapsed="sidebarCollapsed"
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
