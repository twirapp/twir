<template>
  <div
    v-if="user"
    class="grid grid-flow-col grid-cols-[220px_1fr] h-screen max-h-screen overflow-hidden"
  >
    <aside class="flex flex-col border-r border-black-25 p-2 justify-between overflow-hidden">
      <ul class="inline-flex flex-col w-full">
        <li v-for="(item, key) in appMenu" :key="key">
          <router-link :to="item.path" class="app__nav-menu-item" activeClass="active">
            <TswIcon :name="appMenuIcons[key]" :width="20" :height="20" />
            {{ menuTranslation[key] }}
          </router-link>
        </li>
      </ul>
      <div class="inline-grid gap-y-[2px] py-1">
        <a
          class="text-xs text-gray-70 flex items-center hover:text-white-95 transition-colors"
          href="/"
        >
          <TswIcon
            name="ArrowInCircle"
            class="stroke-gray-60 m-1"
            :rotate="-90"
            :width="20"
            :height="20"
            :strokeWidth="1.25"
          />
          Upgrade plan
        </a>
        <a
          class="text-xs text-gray-70 flex items-center hover:text-white-95 transition-colors"
          href="#"
        >
          <TswIcon
            name="Message"
            class="stroke-gray-60 m-1"
            :width="20"
            :height="20"
            :strokeWidth="1.25"
          />
          Leave feedback
        </a>
      </div>
    </aside>
    <div class="relative bg-black-15">
      <header
        class="
          bg-black-15
          grid grid-flow-col
          gap-x-2
          backdrop-blur-xl
          w-full
          py-[10px]
          px-[14px]
          justify-end
          border-b border-b-black-25
        "
      >
        <button class="p-[6px]">
          <TswIcon name="Bell" class="stroke-gray-70" />
        </button>
        <TswAvatar :src="user.profile_image_url" />
      </header>
      <main class="relative px-9">
        <router-view v-slot="{ Component }">
          <Suspense>
            <component :is="Component" />
          </Suspense>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { TswIcon, TswAvatar } from '@tsuwari/ui-components';

import { appMenu, appMenuIcons } from './router.js';

import { useUserProfile } from '@/services/auth';
import { useTranslation } from '@/services/locale';

const { tm } = useTranslation<'app'>();
const { data: user } = useUserProfile();

const menuTranslation = tm('pages');
</script>

<style lang="postcss" scoped>
.app__nav-menu-item {
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
