<template>
  <header :ref="(h) => (headerStore.set(h as HTMLElement))">
    <div class="flex container py-3 items-center justify-between">
      <div class="flex-1 flex">
        <div class="mr-auto flex items-center justify-between max-lg:w-full">
          <a class="inline-grid items-center grid-flow-col gap-x-[10px] p-2 max-lg:p-1" href="#">
            <img :src="Logo" alt="Tsuwari logo" />
            <span class="font-medium text-xl">Tsuwari</span>
          </a>
          <BurgerMenuButton />
          <ClientOnly>
            <MobileMenu />
          </ClientOnly>
        </div>
      </div>
      <NavMenu
        menuItemClass="header-nav-link"
        menuClass="inline-grid grid-flow-col gap-x-2"
        class="max-lg:hidden"
      />
      <div class="flex-1 flex max-lg:bg-red-60 max-lg:hidden">
        <div class="inline-grid grid-flow-col gap-x-3 items-center ml-auto">
          <LangSelect @change="setLandingLocale" />
          <a href="#" class="login-btn">{{ t('buttons.login') }}</a>
        </div>
      </div>
    </div>
    <ClientOnly>
      <div
        :class="{
          'header-bottom-line': true,
          active: windowY > headerHeight,
        }"
      ></div>
    </ClientOnly>
  </header>
</template>

<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { useWindowScroll } from '@vueuse/core';

import Logo from '@/assets/NewLogo.svg';
import ClientOnly from '@/components/ClientOnly.vue';
import BurgerMenuButton from '@/components/landing/layout/BurgerMenuButton.vue';
import MobileMenu from '@/components/landing/layout/MobileMenu.vue';
import NavMenu from '@/components/landing/layout/NavMenu.vue';
import LangSelect from '@/components/LangSelect/LangSelect.vue';
import useLandingLocale from '@/hooks/useLandingLocale';
import useTranslation from '@/hooks/useTranslation';
import { headerStore, headerHeightStore } from '@/stores/landing/header.js';

const setLandingLocale = useLandingLocale();

const headerHeight = useStore(headerHeightStore);

const t = useTranslation<'landing'>();

const { y: windowY } = useWindowScroll();
</script>

<style lang="postcss">
header {
  @apply sticky
    w-full
    left-0
    right-0
    top-0
    mx-auto
    z-20
    bg-black-10 bg-opacity-70
    backdrop-blur-sm backdrop-saturate-[180%];

  .container {
    @apply max-w-[1200px] px-6 max-lg:px-4;
  }
}

.header-bottom-line {
  @apply w-full h-[1px] bg-black-17 transition-colors bg-opacity-0 absolute bottom-0;

  &.active {
    @apply bg-opacity-100;
  }
}

.header-nav-link {
  @apply leading-tight
    cursor-pointer
    px-3
    py-[10px]
    text-gray-70
    hover:text-white-100;

  transition: color 0.15s theme('transitionTimingFunction.in-out');
}

.login-btn {
  @apply inline-flex
    bg-purple-60
    px-[13px]
    py-[7px]
    rounded-md
    leading-tight
    hover:bg-opacity-20
    hover:border-opacity-50
    hover:text-purple-95
    border-2 border-opacity-0 border-purple-70
    transition-colors;
}
</style>
