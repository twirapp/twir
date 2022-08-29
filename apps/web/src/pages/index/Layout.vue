<template>
  <Header :menuItems="navMenuItems" />
  <main>
    <slot />
  </main>
  <Footer :menuItems="navMenuItems" :socials="socials" />
</template>

<script lang="ts" setup>
import Footer from '@/components/sections/Footer.vue';
import Header from '@/components/sections/Header.vue';
import { socials, navMenuItems } from '@/data';
</script>