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
import { isClient } from '@vueuse/core';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { type NavMenuLocale, navMenuHrefs } from '@/data/landing/navMenu.js';
import { useLandingHeaderHeight } from '@/services/landing-menu/header.js';

const props =
  defineProps<{
    menuClass: string;
    menuItemClass: string;
    menuItemClickHandler?: () => any;
  }>();

const { tm } = useI18n();
const headerHeight = useLandingHeaderHeight();

const menuItems = computed(() => tm('navMenu') as NavMenuLocale[]);

const scrollToSection = (e: Event) => {
  if (!isClient) return;

  const sectionId = (e.target as HTMLLinkElement).dataset.section as string;
  const section = document.getElementById(sectionId);
  if (!section) {
    return console.error(`Section "${sectionId}" is not found`);
  }

  const sectionY = window.scrollY - headerHeight.value + section.getBoundingClientRect().top;
  window.scrollTo({
    top: sectionY,
    behavior: 'smooth',
  });

  if (props.menuItemClickHandler) {
    props.menuItemClickHandler();
  }
};
</script>
