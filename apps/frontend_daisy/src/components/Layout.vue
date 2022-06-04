<script lang="ts" setup>
import Footer from './Footer.vue';
import Header from './Header.vue';
import Sidebar from './Sidebar.vue';
</script>

<template>
  <main>
    <div style="font-family: 'Inter', sans-serif">
      <div class="flex flex-row min-h-screen">
        <Sidebar />

        <main class="main flex flex-col flex-grow">
          <Header />

          <div class="flex flex-col flex-grow">
            <slot />
          </div>

          <Footer />
        </main>
      </div>
    </div>
  </main>
</template>
