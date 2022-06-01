<script lang="ts" setup>
import Footer from './Footer.vue';
import NavBar from './NavBar.vue';
import Sidebar from './Sidebar.vue';
</script>

<template>
  <main>
    <div class="flex flex-row min-h-screen text-white">
      <Sidebar />

      <div class="main flex flex-col flex-grow">
        <NavBar />

        <div class="flex flex-col flex-grow m-3">
          <slot />
        </div>

        <Footer />
      </div>
    </div>
  </main>
</template>
