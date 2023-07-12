<script setup lang="ts">
import {
	darkTheme, NLayout, NLayoutHeader, NLayoutContent, NLayoutSider, NConfigProvider, NButton } from 'naive-ui';
import { computed, ref } from 'vue';
import { RouterView } from 'vue-router';

import { useTheme } from './hooks/index.js';
import Header from './layout/header.vue';
import Sidebar from './layout/sidebar.vue';

const localStorageTheme = useTheme();
const theme = computed(() => localStorageTheme.value === 'dark' ? darkTheme : null);

const sidebarCollapsed = ref(false);
function toggleSidebar() {
	sidebarCollapsed.value = !sidebarCollapsed.value;
}
</script>

<template>
  <n-config-provider :theme="theme as any" style="height: 100%">
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
        <n-layout-content content-style="padding: 24px">
          <router-view v-slot="{ Component, route }">
            <transition :name="route.meta.transition || 'fade'">
              <component :is="Component" />
            </transition>
          </router-view>
        </n-layout-content>
      </n-layout>
    </n-layout>
  </n-config-provider>
</template>
