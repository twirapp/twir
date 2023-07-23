<template>
  <div v-if="user" class="flex flex-col h-screen max-h-screen overflow-hidden">
    <header
      class="
        flex
        justify-between
        bg-black-10
        backdrop-blur-xl
        w-full
        py-[10px]
        px-4
        border-b border-b-black-25
      "
    >
      <a class="inline-grid items-center grid-flow-col gap-x-[10px] p-[3px]" href="#">
        <div class="h-[30px] w-[30px]" :style="{ backgroundImage: cssURL(TsuwariLogo) }" />
        <span class="font-medium text-xl">Tsuwari</span>
      </a>
      <div class="grid grid-flow-col gap-x-2">
        <button class="p-[6px]">
          <TswIcon name="Bell" class="stroke-gray-70" />
        </button>
        <TswDropdown>
          <template #button="{ onClick }">
            <TswAvatar :src="user.avatar" @click="onClick" />
          </template>
          <template #menu>
            <div class="profile-dropdown">
              <!-- <li v-for="item in user.dashboards" :key="item.id">
                {{ JSON.stringify(item) }}
              </li> -->
              <button
                class="grid grid-flow-col gap-2 px-3 py-2 items-center hover:bg-black-25"
                @click="logoutAndGoToLanding"
              >
                <TswIcon name="Logout" class="stroke-gray-70" :width="20" :height="20" />
                <span class="pr-2 text-white-95 text-sm">Logout</span>
              </button>
            </div>
          </template>
        </TswDropdown>
      </div>
    </header>
    <div class="grid grid-flow-col grid-cols-[240px_1fr] flex-1">
      <aside class="flex flex-col border-r border-black-25 p-3 justify-between overflow-hidden">
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
      <main class="relative px-9 bg-black-15">
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
import { TswIcon, TswAvatar, cssURL, TswDropdown } from '@twir/ui-components';

import { appMenu, appMenuIcons } from './router.js';

import TsuwariLogo from '@/assets/brand/TsuwariInCircle.svg';
import { useUserProfile, logout, redirectToLanding } from '@/services/auth';
import { useTranslation } from '@/services/locale';

const { tm } = useTranslation<'app'>();
const { data: user } = useUserProfile();

const menuTranslation = tm('pages');

const logoutAndGoToLanding = async () => {
  const isLoggedOut = await logout();
  if (isLoggedOut) {
    return redirectToLanding();
  }
  console.error('Cannot logout!');
};
</script>

<style lang="postcss" scoped>
.profile-dropdown {
  @apply bg-[#2D2D2D] rounded-md overflow-hidden border border-[#414141];

  box-shadow: 0px 2px 36px rgba(0, 0, 0, 0.08), 0px 2px 4px rgba(0, 0, 0, 0.08),
    0px 10px 20px rgba(0, 0, 0, 0.15);
}
.app__nav-menu-item {
  @apply inline-grid
    grid-flow-col
    items-center
    gap-x-2
    p-[6px]
    text-sm
    w-full
    justify-start
    border-2
    border-black-10
    text-[#C5C5C8]
    rounded-md;

  transition-property: background-color border-color color;
  transition-duration: 200ms;
  transition-timing-function: theme('transitionTimingFunction.out');

  &:hover {
    background: theme('colors.black.25');
  }

  & > svg {
    @apply stroke-gray-70 m-[2px];
    transition: stroke 200ms theme('transitionTimingFunction.out');
  }

  &.active {
    @apply hover:bg-purple-55 border-[#816ef2] text-white-100 p-[6px];
    box-shadow: 0px 2px 8px rgba(100, 78, 232, 0.1);
    background-color: #5845c9;

    & > svg {
      @apply stroke-white-100;
    }
  }
}
</style>
