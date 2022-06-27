<script lang="ts" setup>
import Footer from './Footer.vue';
import NavBar from './NavBar.vue';
import Sidebar from './Sidebar.vue';
</script>

<template>
  <main>
    <div class="flex h-screen text-white">
      <Sidebar />

      <div class="flex flex-col flex-grow main">
        <NavBar />

        <div
          class="flex flex-col flex-grow max-h-screen overflow-y-auto"
        >
          <slot />
        </div>

        <Footer />
      </div>
    </div>
  </main>
</template>
