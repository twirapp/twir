<script lang="ts" setup>
import { useRouter } from 'vue-router';

import Footer from './Footer.vue';
import NavBar from './NavBar.vue';
import Sidebar from './Sidebar.vue';

const router = useRouter();
</script>

<template>
  <main>
    <div class="flex flex-row h-screen text-white">
      <Sidebar />

      <div class="main flex flex-col flex-grow">
        <NavBar />

        <div
          class="flex flex-col flex-grow"
          :class="{'m-3': router.currentRoute.value.name !== 'Commands'}"
        >
          <slot />
        </div>

        <Footer />
      </div>
    </div>
  </main>
</template>
