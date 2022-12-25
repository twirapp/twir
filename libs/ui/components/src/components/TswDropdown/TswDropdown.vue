<template>
  <div class="relative inline-block">
    <div ref="button" class="flex">
      <slot name="button" :isOpen="isOpen" :onClick="toggleMenu" />
    </div>
    <Transition name="tsw-dropdown">
      <div v-if="isOpen" ref="menu" class="absolute right-0 mt-2" role="menu">
        <slot name="menu" />
      </div>
    </Transition>
  </div>
</template>

<script lang="ts" setup>
import { onClickOutside } from '@vueuse/core';
import { ref } from 'vue';

const isOpen = ref(false);
const menu = ref<HTMLElement | null>(null);
const button = ref<HTMLElement | null>(null);

onClickOutside(menu as any, (event) => {
  if (!button.value || !event.target) return;

  if (!button.value.contains(event.target as Node)) {
    isOpen.value = false;
  }
});

const toggleMenu = () => (isOpen.value = !isOpen.value);
</script>

<style lang="postcss">
.tsw-dropdown-enter-active,
.tsw-dropdown-leave-active {
  transition: transform 0.18s theme('transitionTimingFunction.DEFAULT'),
    opacity 0.18s theme('transitionTimingFunction.DEFAULT');
}

.tsw-dropdown-enter-from,
.tsw-dropdown-leave-to {
  transform: translateY(-10px) scale(0.9);
  transform-origin: top top;
  opacity: 0;
}
</style>
