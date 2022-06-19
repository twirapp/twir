<script lang="ts" setup>
import { useRouter } from 'vue-router';

import Footer from './Footer.vue';
import NavBar from './NavBar.vue';
import Sidebar from './Sidebar.vue';

const router = useRouter();
</script>

<template>
  <main>
    <div class="flex h-screen text-white">
      <Sidebar />

      <div class="main flex flex-col flex-grow">
        <NavBar />

        <div
          class="flex max-h-screen overflow-y-auto flex-col flex-grow scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-600 scrollbar-track-gray-900"
        >
          <slot />
        </div>

        <Footer />
      </div>
    </div>
  </main>
</template>
