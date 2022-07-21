<script lang="ts" setup>
import Verified from '../../assets/icons/verified.svg?component';

import LandingButtons from '@/components/LandingButtons.vue';
</script>

<template>
  <div class="flex h-screen text-center">
    <div class="m-auto space-y-16 text-white">
      <div>
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

      <div class="flex justify-center">
        <a
          href="#install"
        >
          <div class="animate-bounce bg-[#191919] flex h-10 hover:bg-[#202020] items-center justify-center p-2 ring-1 ring-slate-200/20 rounded-full w-10">
            <Verified />
          </div>
        </a>
      </div>
    </div>
  </div>
</template>