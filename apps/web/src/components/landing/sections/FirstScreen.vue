<template>
  <section class="relative overflow-hidden">
    <AsyncBlobs v-if="isBrowser" />
    <div class="container">
      <div class="inline-flex flex-col items-center w-full pt-28 pb-40">
        <h1 class="font-bold text-7xl text-center">
          The
          <span class="gradient-title">perfect bot</span><br />
          for your stream
        </h1>
        <p class="text-center leading-[145%] text-gray-70 text-[22px] pt-6">
          Created by streamers. Made for streamers.<br />
          Used by streamers. For streamers with love
        </p>
        <div class="inline-grid grid-flow-col gap-x-3 w-[320px] mt-12">
          <a
            class="
              inline-flex
              text-white-95
              bg-black-10
              border border-black-25
              px-5
              py-3
              rounded-full
              justify-center
              leading-tight
            "
          >
            Learn more
          </a>
          <a class="inline-flex bg-purple-60 px-5 py-3 rounded-full justify-center leading-tight">
            Start for free
          </a>
        </div>
      </div>
      <img class="absolute -bottom-[438px] -z-10 mx-auto left-0 right-0" src="@/assets/Waves.svg" />
      <img src="@/assets/Rhombus.svg" class="absolute top-[15%] right-[24%]" />
      <img class="absolute top-[41%] left-[22%]" src="@/assets/Lighting.svg" />
      <img src="@/assets/SmileBot.svg" class="absolute right-[216px] bottom-[200px]" />
    </div>
  </section>
</template>

<script lang="ts" setup>
import { defineAsyncComponent } from 'vue';

import { useBrowser } from '@/hooks/useBrowser.js';

const { isBrowser } = useBrowser();

const AsyncBlobs = defineAsyncComponent(async () => {
  return (await import('@/components/landing/FirstScreenBlobs.vue')).default;
});
</script>

<style lang="postcss">
.gradient-title {
  background: linear-gradient(258.67deg, #d34bf4 4.22%, #905bff 73.01%);

  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  line-height: 1.15;
}
</style>
