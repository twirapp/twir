<template>
  <Header :menuItems="tm('navMenu')" />
  <main>
    <slot />
  </main>
  <Footer :menuItems="tm('navMenu')" :socials="socials" />
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import Footer from '@/components/sections/Footer.vue';
import Header from '@/components/sections/Header.vue';
import { socials } from '@/data';

const { tm } = useI18n();
</script>