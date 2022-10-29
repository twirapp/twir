<template>
  <div class="relative inline-flex">
    <button ref="select" :class="`select ${isOpen ? 'open' : ''}`" @click="isOpen = !isOpen">
      {{ pageContext.locale.toUpperCase() }}
      <div class="icon">
        <TswIcon name="ArrowTriangleMedium" :rotate="90" />
      </div>
    </button>
    <Transition>
      <div v-if="isOpen" ref="dropdownRef" class="dropdown">
        <LangSelectOption
          v-for="lang in languages"
          :key="lang.locale"
          :isActive="lang.locale === pageContext.locale"
          :locale="lang.locale"
          :name="lang.name"
          @change="(l) => emit('change', l)"
        />
      </div>
    </Transition>
  </div>
</template>

<script lang="ts" setup>
import { TswIcon } from '@tsuwari/ui-components';
import { onClickOutside } from '@vueuse/core';
import { ref } from 'vue';

import LangSelectOption from './LangSelectOption.vue';

import type { Locale } from '@/locales';
import { languages } from '@/locales';
import { usePageContext } from '@/utils/pageContext.js';

const emit = defineEmits<{ (event: 'change', locale: Locale): void }>();

const dropdownRef = ref<HTMLElement | null>(null);
const select = ref<HTMLElement | null>(null);

const pageContext = usePageContext();

onClickOutside(dropdownRef, (event) => {
  if (!select.value) return;

  if (!select.value.contains(event.target as HTMLElement)) {
    isOpen.value = false;
  }
});

const isOpen = ref(false);
</script>

<style scoped lang="postcss">
.v-enter-active,
.v-leave-active {
  transition: transform 0.25s theme('transitionTimingFunction.DEFAULT'),
    opacity 0.25s theme('transitionTimingFunction.DEFAULT');
}

.v-enter-from,
.v-leave-to {
  transform: translateY(10px);
  opacity: 0;
}

.select {
  @apply inline-grid
    grid-flow-col
    gap-x-[6px]
    items-center
    py-[10px]
    px-3
    leading-tight;

  &:hover > .icon {
    transform: translateY(2px);
  }

  & > .icon {
    @apply stroke-gray-60 transition-transform;
  }

  &.open > .icon {
    transform: translateY(2px);
  }
}

.dropdown {
  @apply inline-flex
    flex-col
    p-[5px]
    bg-black-15
    border border-black-25
    rounded-md
    min-w-[100px]
    absolute
    right-0
    top-10
    z-10;
}
</style>
