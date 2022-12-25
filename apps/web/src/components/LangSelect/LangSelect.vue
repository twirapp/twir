<template>
  <TswDropdown>
    <template #button="{ isOpen, onClick }">
      <button :class="`select ${isOpen ? 'open' : ''}`" @click="onClick">
        {{ pageContext.locale.toUpperCase() }}
        <div class="icon">
          <TswIcon name="ArrowTriangleMedium" :rotate="90" />
        </div>
      </button>
    </template>
    <template #menu>
      <div class="dropdown">
        <LangSelectOption
          v-for="lang in languages"
          :key="lang.locale"
          :locale="lang.locale"
          :name="lang.name"
          :isActive="pageContext.locale === lang.locale"
          @change="(l) => emit('change', l)"
        />
      </div>
    </template>
  </TswDropdown>
</template>

<script lang="ts" setup>
import { TswIcon, TswDropdown } from '@tsuwari/ui-components';

import LangSelectOption from './LangSelectOption.vue';

import type { Locale } from '@/locales';
import { languages } from '@/locales';
import { usePageContext } from '@/utils/pageContext.js';

const emit = defineEmits<{ (event: 'change', locale: Locale): void }>();

const pageContext = usePageContext();
</script>

<style scoped lang="postcss">
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
    p-1
    bg-black-15
    border border-black-25
    rounded-[8px]
    min-w-[100px]
    z-10
    top-0;
}
</style>
