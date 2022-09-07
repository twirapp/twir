<template>
  <header :ref="(h) => (headerStore.set(h as HTMLElement))">
    <div class="flex container py-3 items-center justify-between">
      <div class="flex-1 flex">
        <div class="mr-auto">
          <a class="inline-grid items-center grid-flow-col gap-x-[10px] p-2" href="#">
            <img :src="Logo" alt="Tsuwari logo" />
            <span class="font-medium text-xl">Tsuwari</span>
          </a>
        </div>
      </div>
      <NavMenu
        menuItemClass="header-nav-link"
        menuClass="inline-grid grid-flow-col gap-x-2"
        :menuItems="menuItems"
      />
      <div class="flex-1 flex">
        <div class="inline-grid grid-flow-col gap-x-3 items-center ml-auto">
          <LangSelect @change="setLocale" />
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
import { navigate } from 'vite-plugin-ssr/client/router';
import { useI18n } from 'vue-i18n';

import Logo from '@/assets/NewLogo.svg';
import ClientOnly from '@/components/ClientOnly.vue';
import NavMenu from '@/components/landing/layout/NavMenu.vue';
import LangSelect from '@/components/LangSelect/LangSelect.vue';
import { headerStore, headerHeightStore } from '@/stores/landing/header.js';
import type { Locale } from '@/types/locale.js';
import type { NavMenuLocale } from '@/types/navMenu';
import { loadLocaleMessages, useTranslation } from '@/utils/locales.js';

defineProps<{ menuItems: NavMenuLocale[] }>();

const headerHeight = useStore(headerHeightStore);

const t = useTranslation<'landing'>();

const { setLocaleMessage, locale: i18nLocale } = useI18n();

async function setLocale(locale: Locale) {
  const messages = await loadLocaleMessages('landing', locale);

  setLocaleMessage<any>(locale, messages);
  i18nLocale.value = locale;

  navigate(`/${locale}`, { keepScrollPosition: true });
}

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
    bg-black-10 bg-opacity-80
    backdrop-blur-sm backdrop-saturate-[180%];
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
