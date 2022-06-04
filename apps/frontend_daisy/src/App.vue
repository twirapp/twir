<script setup lang="ts">
import { useTitle } from '@vueuse/core';
import { onMounted } from 'vue';
import { RouterView, useRouter } from 'vue-router';

import Layout from '@/components/Layout.vue';


const router = useRouter();

const title = useTitle();
title.value = 'Tsuwari - Main page';

onMounted(async () => {
  /* const accessToken = localStorage.getItem("accessToken");
  const refreshToken = localStorage.getItem("refreshToken");
  if (accessToken && refreshToken) {
    fetchAndSetUser();
  } */
});
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
