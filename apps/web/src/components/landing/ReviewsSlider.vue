<template>
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
    @mouseenter="stopSlider"
    @mouseleave="startSlider"
  >
    <SwiperSlide v-for="item in reviews" :key="item.id" style="flex-shrink: 1">
      <ReviewCard
        class="slider-review-card"
        :username="item.username"
        :comment="item.comment"
        :rating="item.rating"
        :avatarUrl="item.avatarUrl"
      />
    </SwiperSlide>
  </Swiper>
</template>

<script lang="ts" setup>
import { Autoplay, type Swiper as ISwiper } from 'swiper';
import { Swiper, SwiperSlide } from 'swiper/vue';
import { ref } from 'vue';

import 'swiper/css';

import ReviewCard from '@/components/landing/ReviewCard.vue';
import type { Review } from '@/data/landing/reviews.js';

defineProps<{
  reviews: Review[]
}>();

const slider = ref<ISwiper | null>(null);

const setSwiper = (swiper: ISwiper) => {
  slider.value = swiper;
};

const stopSlider = () => {
  if (slider.value) {
    slider.value.autoplay.stop();
  }
};

const startSlider = () => {
  if (slider.value) {
    slider.value.autoplay.start();
  }
};

const modules = [Autoplay];
</script>

<style lang="postcss">
.slider-review-card {
  width: 380px;

  @media screen and (max-width: 565.98px) {
    width: calc(100vw - 24px * 2);
  }
}
</style>
