<template>
  <nav>
    <ul :class="menuClass">
      <li v-for="item in menuItems" :key="item.id" class="inline-flex">
        <button
          :data-section="LandingSection[item.id]"
          :class="menuItemClass"
          @click.prevent="scrollToSection"
        >
          {{ item.name }}
        </button>
      </li>
    </ul>
  </nav>
</template>

<script lang="ts" setup>
import { isClient } from '@vueuse/core';
import { computed } from 'vue';

import { LandingSection } from '@/data/landing/sections.js';
import { scrollToLandingSection } from '@/services/landing';
import { useTranslation } from '@/services/locale/hooks.js';

const props =
  defineProps<{
    menuClass: string;
    menuItemClass: string;
    menuItemClickHandler?: () => any;
  }>();

const { tm } = useTranslation<'landing'>();
const menuItems = computed(() => tm('navMenu'));

const scrollToSection = (e: Event) => {
  if (!isClient) return;

  scrollToLandingSection(e.target as HTMLElement);

  if (props.menuItemClickHandler) {
    props.menuItemClickHandler();
  }
};
</script>
