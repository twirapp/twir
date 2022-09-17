<template>
  <Transition>
    <div
      v-if="menuState"
      class="mobile-menu"
      :style="{
        top: cssPX(headerHeight),
        height: cssPX(windowHeight - headerHeight),
      }"
    >
      <a
        href="#"
        class="
          inline-flex
          bg-purple-60
          px-4
          py-[10px]
          rounded-md
          w-full
          text-center
          justify-center
          hover:bg-purple-50
          transition-colors
        "
      >
        {{ t('buttons.login') }}
      </a>
      <div class="flex flex-col mt-5">
        <NavMenu
          menuClass="mobile-nav-menu"
          menuItemClass="mobile-nav-menu-item"
          :menuItemClickHandler="closeMenu"
        />
        <div class="mt-2 w-full flex justify-end">
          <LangSelect @change="changeLangAndCloseMenu" />
        </div>
      </div>
    </div>
  </Transition>
</template>

<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { useWindowSize } from '@vueuse/core';
import { onUnmounted } from 'vue';

import NavMenu from '@/components/landing/layout/NavMenu.vue';
import LangSelect from '@/components/LangSelect/LangSelect.vue';
import useLandingLocale from '@/hooks/useLandingLocale.js';
import useTranslation from '@/hooks/useTranslation.js';
import type { Locale } from '@/locales';
import { headerHeightStore, menuStateStore } from '@/stores/landing/header.js';
import { cssPX } from '@/utils/css';

const menuState = useStore(menuStateStore);
const headerHeight = useStore(headerHeightStore);

const setLandingLocale = useLandingLocale();

const closeMenu = () => menuStateStore.set(false);

const changeLangAndCloseMenu = (locale: Locale) => {
  setLandingLocale(locale);
  // ??? Do I need to close the menu when changing the language?
  closeMenu();
};

const removeListener = menuStateStore.listen((menuState) => {
  if (menuState) {
    document.body.classList.add('overflow-hidden');
  } else {
    document.body.classList.remove('overflow-hidden');
  }
});

const t = useTranslation<'landing'>();

const { height: windowHeight } = useWindowSize();

onUnmounted(() => {
  removeListener();
});
</script>

<style lang="postcss">
.mobile-nav-menu {
  @apply flex w-full flex-col;

  & > :not(:last-child) {
    @apply border-b border-black-25;
  }
}

.mobile-nav-menu-item {
  @apply flex w-full py-[15px] leading-tight;
}
</style>

<style lang="postcss" scoped>
.mobile-menu {
  @apply block
    fixed
    w-full
    left-0
    right-0
    bottom-0
    max-w-[100vw]
    z-50
    bg-black-10
    p-5;
}

.v-enter-active,
.v-leave-active {
  transition: opacity 0.3s theme('transitionTimingFunction.out'),
    transform 0.3s theme('transitionTimingFunction.out');
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
  transform: translateY(12px);
}
</style>
