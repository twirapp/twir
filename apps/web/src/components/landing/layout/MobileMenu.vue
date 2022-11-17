<template>
  <Transition>
    <div v-if="menuState" class="mobile-menu" :style="menuStyles">
      <div v-if="!user" class="px-5 mb-3">
        <TswButton
          class="w-full text-center justify-center"
          :text="t('buttons.login')"
          @click="redirectToLogin"
        />
      </div>
      <div class="flex flex-col">
        <NavMenu
          menuClass="mobile-nav-menu"
          menuItemClass="mobile-nav-menu-item"
          :menuItemClickHandler="closeMenu"
        />
        <div class="mt-2 w-full flex justify-end px-5">
          <LangSelect @change="changeLangAndCloseMenu" />
        </div>
      </div>
    </div>
  </Transition>
</template>

<script lang="ts" setup>
import { cssPX, TswButton } from '@tsuwari/ui-components';
import { useWindowSize } from '@vueuse/core';
import { computed, StyleValue } from 'vue';

import NavMenu from '@/components/landing/layout/NavMenu.vue';
import LangSelect from '@/components/LangSelect/LangSelect.vue';
import type { Locale } from '@/locales';
import { redirectToLogin } from '@/services/auth';
import { useUserProfile } from '@/services/auth';
import { useLandingHeaderHeight, useLandingMenuState } from '@/services/landing-menu';
import { useLandingLocale, useTranslation } from '@/services/locale';

const { data: user } = useUserProfile();

const { setLandingLocale } = useLandingLocale();
const { t } = useTranslation<'landing'>();

const { menuState, closeMenu } = useLandingMenuState();
const headerHeight = useLandingHeaderHeight();
const { height: windowHeight } = useWindowSize();

const menuStyles = computed<StyleValue>(() => ({
  top: cssPX(headerHeight.value),
  height: cssPX(windowHeight.value - headerHeight.value),
}));

const changeLangAndCloseMenu = (locale: Locale) => {
  setLandingLocale(locale);
  closeMenu();
};
</script>

<style lang="postcss">
.mobile-nav-menu {
  @apply flex w-full flex-col rounded overflow-hidden;

  & > :not(:last-child) {
    @apply border-b border-black-25;
  }
}

.mobile-nav-menu-item {
  @apply flex w-full py-[15px] leading-tight hover:bg-black-15 transition-colors justify-center;
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
    border-t
    border-black-25
    bg-black-10
    pt-5;
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
