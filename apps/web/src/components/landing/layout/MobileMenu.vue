<template>
  <Transition>
    <div
      v-if="menuState"
      class="mobile-menu bg-black-10 p-5"
      :style="{
        top: `${headerHeight}px`,
        height: `${windowHeight - headerHeight}px`,
      }"
    >
      <a
        href="#"
        class="inline-flex bg-purple-60 px-4 py-[10px] rounded-md w-full text-center justify-center"
      >
        {{ t('buttons.login') }}
      </a>
      <LangSelect />
    </div>
  </Transition>
</template>

<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { useWindowSize } from '@vueuse/core';
import { watch } from 'vue';

import LangSelect from '@/components/LangSelect/LangSelect.vue';
import { headerHeightStore } from '@/stores/landing/header.js';
import { useTranslation } from '@/utils/locales.js';

const props =
  defineProps<{
    menuState: boolean;
  }>();

watch(
  () => props.menuState,
  (menuState) => {
    if (menuState) {
      document.body.classList.add('overflow-hidden');
    } else {
      document.body.classList.remove('overflow-hidden');
    }
  },
);

const t = useTranslation<'landing'>();

const headerHeight = useStore(headerHeightStore);
const { height: windowHeight } = useWindowSize();
</script>

<style lang="postcss" scoped>
.mobile-menu {
  @apply block
    fixed
    w-full
    left-0
    right-0
    bottom-0
    max-w-[100vw]
    z-50;
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
