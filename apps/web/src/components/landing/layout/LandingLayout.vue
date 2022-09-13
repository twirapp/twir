<template>
  <div class="relative min-h-screen">
    <Header />
    <main class="overflow-x-hidden">
      <slot />
    </main>
  </div>
  <Footer />
</template>

<script lang="ts" setup>
import Footer from '@/components/landing/layout/Footer.vue';
import Header from '@/components/landing/layout/Header.vue';
</script>
