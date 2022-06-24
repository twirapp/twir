<script lang="ts" setup>
import Footer from './Footer.vue';
import NavBar from './NavBar.vue';
import Sidebar from './Sidebar.vue';
</script>

<template>
  <main>
    <div class="flex h-screen text-white">
      <Sidebar />

      <div class="main flex flex-col flex-grow">
        <NavBar />

        <div
          class="flex max-h-screen flex-col flex-grow overflow-y-auto"
        >
          <slot />
        </div>

        <Footer />
      </div>
    </div>
  </main>
</template>
