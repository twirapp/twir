<template>
  <router-view v-slot="{ Component }">
    <Suspense>
      <component :is="Component" />
    </Suspense>
  </router-view>
  <ul class="inline-flex flex-col min-w-[200px]">
    <li v-for="(item, key) in appMenu" :key="key">
      <router-link :to="item.path" class="nav-menu-item" activeClass="active">
        <TswIcon :name="appMenuIcons[key]" :width="20" :height="20" />
        {{ menuTranslation[key] }}
      </router-link>
    </li>
  </ul>
</template>

<script lang="ts" setup>
import { TswIcon } from '@tsuwari/ui-components';

import { appMenu, appMenuIcons } from './router.js';

import { useTranslation } from '@/services/locale';

const { tm } = useTranslation<'app'>();

const menuTranslation = tm('pages');
</script>

<style lang="postcss">
.nav-menu-item {
  @apply inline-grid
    grid-flow-col
    items-center
    gap-x-2
    p-2
    text-sm
    w-full
    justify-start
    text-white-95
    rounded-md
    hover:bg-black-25;

  & > svg {
    @apply stroke-gray-70 m-[2px];
  }

  &.active {
    @apply bg-purple-60 hover:bg-purple-55 text-white-100;

    & > svg {
      @apply stroke-white-100;
    }
  }
}
</style>
