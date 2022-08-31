<template>
  <section class="bg-black-10 relative overflow-hidden z-0">
    <div class="container relative max-w-[980px] mt-20 mb-16">
      <img src="@/assets/CurveArrow.svg" class="absolute top-[30px] -left-[90px] -z-[1]" />
      <img
        src="@/assets/BlobPurple.svg"
        class="absolute -top-[600px] -right-[400px] -z-10 -rotate-90"
      />
      <img src="@/assets/MessageCircle.svg" class="absolute -z-[1] -bottom-3 right-28" />
      <h2 class="text-5xl font-semibold max-w-xl leading-[125%]">
        Reviews from streamers and other viewers
      </h2>
    </div>
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
  </section>
</template>

<script lang="ts" setup>
import { Autoplay, type Swiper as ISwiper } from 'swiper';
import { Swiper, SwiperSlide } from 'swiper/vue';
import { ref } from 'vue';

import 'swiper/css';

import ReviewCard from '@/components/landing/ReviewCard.vue';
import type { Review } from '@/types/review';

const slider = ref<ISwiper | null>(null);

const setSwiper = (swiper: ISwiper) => {
  slider.value = swiper;
};

const modules = [Autoplay];

defineProps<{reviews: Review[]}>();
</script>
