<script lang="ts" setup>
import { HelixStreamData } from '@twurple/api';
import { useIntervalFn, useTitle } from '@vueuse/core';
import { intervalToDuration, formatDuration } from 'date-fns';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';

import CmdMenu from '../../assets/buttons/cmd_menu.svg?component';
import LanguageSelector from '../LanguageSelector.vue';
import Notification from '../Notification.vue';
import Profile from '../Profile.vue';
import { publicRoutes, adminRoutes } from './sidebarLinks';

import { useUpdatingData } from '@/functions/useUpdatingData';

const { data: stream, execute } = useUpdatingData<
  (HelixStreamData & { parsedMessages: number }) | null
>(`/v1/channels/{dashboardId}/streams`);

const uptime = ref('');
const currentRoute = useRoute();
const title = useTitle();
const { t } = useI18n();

useIntervalFn(() => {
  execute();
}, 15000);

useIntervalFn(
  () => {
    if (stream.value) {
      uptime.value = formatDuration(
        intervalToDuration({ start: new Date(stream.value.started_at), end: Date.now() }),
      );
    }
  },
  1000,
  { immediate: true },
);
</script>

<template>
  <nav
    class="border-b border-stone-700 flex flex-wrap items-center justify-between py-3 relative shadow text-white w-full">
    <div class="flex flex-wrap items-center justify-between px-6 w-full">
      <div class="flex space-x-2">
        <div class="block sm:hidden">
          <CmdMenu />

          <ul
            class="absolute bg-[#202020] bg-clip-padding border-none dropdown-menu float-left hidden list-none max-h-[85vh] overflow-y-auto px-1 py-2 rounded space-y-0.5 text-left w-full z-1"
            aria-labelledby="sidebarDropdown">
            <li
              v-for="(route, index) in currentRoute.fullPath.includes('/admin')
                ? adminRoutes
                : publicRoutes"
              :key="index"
              class="flex hover:bg-[#393636] items-center px-2 py-2 rounded space-x-2">
              <RouterLink
                :to="route.path"
                class="bg-transparent flex space-x-2 text-white w-full whitespace-nowrap"
                :class="{
                  'bg-neutral-700 rounded py-0.5 pl-0.5': currentRoute.path === route.path,
                }"
                @click="
                  title = `Tsuwari - ${
                    route.name.charAt(0).toUpperCase() + route.name.substring(1)
                  }`
                ">
                <img :src="route.icon" />

                <span>{{ t(`pages.${route.name}.sidebarName`) }}</span>
              </RouterLink>
            </li>
          </ul>
        </div>
        <div v-if="stream && Object.keys(stream).length" class="flex space-x-2">
          <p>
            Viewers: <span class="font-bold">{{ stream.viewer_count }}</span>
          </p>
          <p class="hidden lg:block">
            Category: <span class="font-bold">{{ stream.game_name }}</span>
          </p>
          <p class="hidden md:block">
            Title:
            <span class="font-bold">{{
              stream.title.length >= 20 ? stream.title.slice(0, 20) + '...' : stream.title
            }}</span>
          </p>
        </div>
        <div v-else class="flex space-x-2">
          {{ t('navbar.offline') }}
        </div>
      </div>

      <div class="flex flex-row space-x-3.5">
        <LanguageSelector />

        <!-- <Notification /> -->
        <Profile />
      </div>
    </div>
  </nav>
</template>
