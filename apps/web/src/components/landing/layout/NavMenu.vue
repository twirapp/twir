<template>
  <nav>
    <ul :class="menuClass">
      <li v-for="item in menuItems" :key="item.id" class="inline-flex">
        <button
          :data-section="navMenuHrefs[item.id]"
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
import { useStore } from '@nanostores/vue';

import { navMenuHrefs } from '@/data/index';
import { headerHeightStore } from '@/stores/landing/header.js';
import type { NavMenuLocale } from '@/types/navMenu.js';

defineProps<{ menuItems: NavMenuLocale[]; menuItemClass: string; menuClass: string }>();

const scrollToSection = (e: Event) => {
  if (typeof window === 'undefined') return;

  const headerHeight = useStore(headerHeightStore);

  const sectionId = (e.target as HTMLLinkElement).dataset.section as string;
  const section = document.getElementById(sectionId);
  if (!section) {
    console.error('Section is not founded');
    return;
  }

  const sectionY = window.scrollY - headerHeight.value + section.getBoundingClientRect().top;

  window.scrollTo({
    top: sectionY,
    behavior: 'smooth',
  });
};
</script>