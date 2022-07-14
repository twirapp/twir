<script lang="ts" setup>
import LandingButtons from '@/components/LandingButtons.vue';
import MainNavBar from '@/components/layout/MainNavBar.vue';
</script>

<template>
  <div style="font-family: 'Inter', sans-serif">
    <MainNavBar />

    <div class="flex h-screen text-center">
      <div class="m-auto text-white">
        <h1 class="break-words font-bold mb-5 md:text-7xl outline-none text-5xl">
          The perfect bot for<br> your stream
        </h1>
        <p class="py-6 text-xl">
          Created by streamers. Made
          <span style="color: #772CE8">for streamers</span>.<br>Used
          <span style="color: #772CE8">by streamers</span>. For streamers with love.
        </p>
   
        <LandingButtons size="large" />
      </div>
    </div>
  </div>
</template>
