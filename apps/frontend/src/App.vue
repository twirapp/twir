<script setup lang="ts">
import { useTitle } from '@vueuse/core';
import { RouterView, useRouter } from 'vue-router';

import Layout from '@/components/layout/Layout.vue';

const router = useRouter();

const title = useTitle();
title.value = 'Tsuwari - Main page';
</script>

<template>
  <Layout v-if="!['/', '/login'].includes(router.currentRoute.value.path)">
    <RouterView v-slot="{ Component }">
      <transition>
        <component :is="Component" />
      </transition>
    </RouterView>
  </Layout>
  <RouterView v-else />
</template>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
