<template>
  <header :ref="setHeaderRef">
    <div class="px-6 max-lg:px-4">
      <div class="container">
        <div class="flex-1 flex">
          <div class="mr-auto flex items-center justify-between max-lg:w-full">
            <a class="inline-grid items-center grid-flow-col gap-x-[10px] p-2 max-lg:p-1" href="#">
              <div class="h-[30px] w-[30px]" :style="{ backgroundImage: cssURL(TsuwariLogo) }" />
              <span class="font-medium text-xl">Tsuwari</span>
            </a>
            <div class="inline-grid grid-flow-col gap-x-3">
              <BurgerMenuButton />
              <ClientOnly>
                <HeaderAuthBlock v-if="!isDesktop" />
              </ClientOnly>
            </div>
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
        <div class="flex-1 flex max-lg:hidden">
          <div class="inline-grid grid-flow-col gap-x-3 items-center ml-auto">
            <LangSelect @change="setLandingLocale" />
            <ClientOnly>
              <HeaderAuthBlock>
                <template #error>
                  <TswButton :text="t('buttons.login')" @click="redirectToLogin" />
                </template>
              </HeaderAuthBlock>
            </ClientOnly>
          </div>
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
import { cssURL, TswButton } from '@tsuwari/ui-components';
import { isClient, useWindowScroll, useWindowSize } from '@vueuse/core';
import { computed } from 'vue';

import TsuwariLogo from '@/assets/brand/TsuwariInCircle.svg';
import ClientOnly from '@/components/ClientOnly.vue';
import BurgerMenuButton from '@/components/landing/layout/BurgerMenuButton.vue';
import HeaderAuthBlock from '@/components/landing/layout/HeaderAuthBlock.vue';
import MobileMenu from '@/components/landing/layout/MobileMenu.vue';
import NavMenu from '@/components/landing/layout/NavMenu.vue';
import LangSelect from '@/components/LangSelect/LangSelect.vue';
import { redirectToLogin } from '@/services/auth';
import { useLandingHeaderHeight, setHeaderRef } from '@/services/landing-menu';
import { useLandingLocale, useTranslation } from '@/services/locale';

const { setLandingLocale } = useLandingLocale();

const headerHeight = useLandingHeaderHeight();

const { y: windowY } = useWindowScroll();
const { width: windowWidth } = useWindowSize();
const { t } = useTranslation();

const MIN_DESKTOP_WIDTH = 996;
const isDesktop = computed(() => windowWidth.value >= MIN_DESKTOP_WIDTH);
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
    @apply flex max-w-[1200px] py-3 items-center justify-between;
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
