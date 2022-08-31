<template>
  <div class="relative min-h-screen">
    <Header :menuItems="tm('navMenu')" />
    <main>
      <slot />
    </main>
  </div>
  <Footer :menuItems="tm('navMenu')" :socials="socials" />
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import Footer from '@/components/landing/layout/Footer.vue';
import Header from '@/components/landing/layout/Header.vue';
import { socials } from '@/data';

const { tm } = useI18n();
</script>
