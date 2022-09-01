<template>
  <section class="bg-black-10 relative overflow-hidden z-0">
    <div class="container relative max-w-[980px] mt-20 mb-16">
      <img src="@/assets/CurveArrow.svg" class="absolute top-[30px] -left-[90px] -z-[1]" />
      <ClientOnly>
        <BlurryBlob
          v-if="isBrowser"
          class="-top-[300px] -right-[140px] w-[490px] h-[276px] rotate-[27deg]"
          color="purple"
        />
      </ClientOnly>
      <img src="@/assets/MessageCircle.svg" class="absolute -z-[1] -bottom-3 right-28" />
      <h2 class="text-5xl font-semibold max-w-xl leading-[125%]">
        Reviews from streamers and other viewers
      </h2>
    </div>
    <ClientOnly>
      <Swiper
        :space-between="24"
        :autoplay="{
          delay: 1,
          disableOnInteraction: false,
        }"
        :speed="2000"
        slidesPerView="auto"
        :centeredSlides="true"
        :centeredSlidesBounds="true"
        :loop="true"
        :grabCursor="true"
        :modules="modules"
        class="mb-20"
        @swiper="setSwiper"
        @mouseenter="slider?.autoplay.stop()"
        @mouseleave="slider?.autoplay.start()"
      >
        <SwiperSlide v-for="item in reviews" :key="item.id" style="width: 380px">
          <ReviewCard
            :username="item.username"
            :comment="item.comment"
            :rating="item.rating"
            :avatarUrl="item.avatarUrl"
          />
        </SwiperSlide>
      </Swiper>
    </ClientOnly>
  </section>
</template>

<script lang="ts" setup>
import { Autoplay, type Swiper as ISwiper } from 'swiper';
import { Swiper, SwiperSlide } from 'swiper/vue';
import { ref } from 'vue';

import 'swiper/css';

import BlurryBlob from '@/components/BlurryBlob.vue';
import ClientOnly from '@/components/ClientOnly';
import ReviewCard from '@/components/landing/ReviewCard.vue';
import { useBrowser } from '@/hooks/useBrowser.js';
import type { Review } from '@/types/review';

defineProps<{reviews: Review[]}>();

const { isBrowser } = useBrowser();

const slider = ref<ISwiper | null>(null);

const setSwiper = (swiper: ISwiper) => {
  slider.value = swiper;
};

const modules = [Autoplay];
</script>
