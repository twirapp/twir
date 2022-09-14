<template>
  <section class="bg-black-15 border-b border-b-black-25">
    <div class="container">
      <div class="flex min-md:w-[740px] mx-auto py-[18px]">
        <ClientOnly :renderClient="renderClient">
          <template #default>
            <Swiper
              tag="ul"
              :space-between="24"
              :speed="1000"
              :slidesPerView="slidesPerView"
              :loop="true"
              :grabCursor="true"
              :autoplay="true"
              :modules="modules"
              class="w-full mx-10 max-sm:mx-0 animate-fadeIn opacity-0"
            >
              <SwiperSlide v-for="item in stats" :key="item.id" class="flex justify-center">
                <StatsItem :item="item" />
              </SwiperSlide>
            </Swiper>
          </template>
          <template #server>
            <ul class="inline-flex gap-x-6 justify-between w-full max-md:opacity-0 overflow-hidden">
              <StatsItem v-for="item in stats" :key="item.id" :item="item" class="w-full flex-1" />
            </ul>
          </template>
        </ClientOnly>
      </div>
    </div>
  </section>
</template>

<script lang="ts" setup>
import { useWindowSize } from '@vueuse/core';
import { Autoplay } from 'swiper';
import { Swiper, SwiperSlide } from 'swiper/vue';
import { computed } from 'vue';

import ClientOnly from '@/components/ClientOnly.vue';
import StatsItem from '@/components/landing/StatsItem.vue';
import { stats } from '@/data/landing/statInfo.js';

import 'swiper/css';

const { width: windowWidth } = useWindowSize();

const renderClient = computed(() => windowWidth.value < 768);

const slidesPerView = computed(() => {
  if (windowWidth.value < 410) {
    return 1;
  }

  if (windowWidth.value < 568) {
    return 2;
  }

  return 3;
});

const modules = [Autoplay];
</script>
