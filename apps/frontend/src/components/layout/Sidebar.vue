<script lang="ts" setup>
import { useTitle } from '@vueuse/core';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';

import Logo from '../../assets/logo.svg?component';
import { publicRoutes, adminRoutes } from './sidebarLinks';

const currentRoute = useRoute();
const { t } = useI18n({
  useScope: 'global',
  inheritLocale: true,
});
const title = useTitle();
</script>

<template>
  <aside
    class="
      border-r border-stone-700
      h-screen
      hidden
      min-w-max
      overflow-auto
      scrollbar
      select-none
      shadow
      sm:block
      w-54
    "
  >
    <div class="flex items-center justify-center py-1 sidebar-header">
      <div>
        <a href="/" class="flex-row font-bold inline-flex items-center mt-5 text-xl">
          <div class="flex items-center space-x-2">
            <Logo />
            <p>Tsuwari</p>
          </div>
        </a>
      </div>
    </div>

    <ul class="mt-3 relative">
      <li
        v-for="(route, index) in currentRoute.fullPath.includes('/admin')
          ? adminRoutes
          : publicRoutes"
        :key="index"
      >
        <RouterLink
          :to="route.path"
          class="
            border-slate-300
            duration-300
            ease-in-out
            flex
            h-12
            hover:bg-[#202122]
            items-center
            mt-1
            overflow-hidden
            px-6
            py-4
            ripple-surface-primary
            rounded
            text-ellipsis text-sm text-white
            transition
            whitespace-nowrap
          "
          :class="{
            'bg-neutral-700': currentRoute.path === route.path,
          }"
          @click="
            title = `Tsuwari - ${route.name.charAt(0).toUpperCase() + route.name.substring(1)}`
          "
        >
          <span v-if="route.icon" class="h-3 mr-3 w-3">
            <img :src="route.icon" />
          </span>

          <span>{{ t(`pages.${route.name}.sidebarName`) }}</span>
        </RouterLink>
      </li>
    </ul>
  </aside>
</template>
