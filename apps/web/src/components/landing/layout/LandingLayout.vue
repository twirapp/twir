<template>
  <div class="relative min-h-screen">
    <Header />
    <main class="overflow-x-hidden">
      <slot />
    </main>
  </div>
  <Footer :socials="socials" />
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import Footer from '@/components/landing/layout/Footer.vue';
import Header from '@/components/landing/layout/Header.vue';
import { socials } from '@/data';
import { navMenuLocaleStore } from '@/stores/landing/navMenu';

const { tm } = useI18n();

navMenuLocaleStore.set(tm('navMenu'));
</script>
