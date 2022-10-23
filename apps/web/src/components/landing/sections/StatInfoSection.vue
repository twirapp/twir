<template>
  <section class="bg-black-15 border-b border-b-black-25">
    <div class="container">
      <div class="flex min-md:w-[740px] mx-auto py-[18px]">
        <ClientOnly :renderClient="isRenderClient">
          <template #default>
            <Carousel
              :speed="1000"
              :autoplay="2000"
              :pauseAutoplayOnHover="true"
              :itemsToShow="slidesPerView"
              snapAlign="start"
              class="flex w-full max-sm:mx-0 animate-fadeIn opacity-0 cursor-grab select-none"
            >
              <Slide v-for="item in stats" :key="item.id" class="flex justify-center">
                <StatsItem :item="item" class="w-full" />
              </Slide>
            </Carousel>
          </template>
          <template #server>
            <div
              class="inline-flex gap-x-6 justify-between w-full max-md:opacity-0 overflow-hidden"
            >
              <StatsItem v-for="item in stats" :key="item.id" :item="item" class="w-full flex-1" />
            </div>
          </template>
        </ClientOnly>
      </div>
    </div>
  </section>
</template>

<script lang="ts" setup>
import { isClient, useWindowSize } from '@vueuse/core';
import { computed } from 'vue';
import { Carousel, Slide } from 'vue3-carousel';

import ClientOnly from '@/components/ClientOnly.vue';
import StatsItem from '@/components/landing/StatsItem.vue';
import { stats } from '@/data/landing/statInfo.js';

import 'vue3-carousel/dist/carousel.css';

const { width: windowWidth } = useWindowSize();

const isRenderClient = computed(() => isClient && windowWidth.value <= 768);

const slidesPerView = computed(() => {
  if (windowWidth.value < 410) {
    return 1;
  } else if (windowWidth.value < 568) {
    return 2;
  } else {
    return 3;
  }
});
</script>
