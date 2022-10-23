<template>
  <Carousel
    :wrapAround="true"
    :autoplay="2500"
    :pauseAutoplayOnHover="true"
    :transition="500"
    :itemsToShow="itemsToShow"
    class="mb-20 cursor-grab"
    :itemsToScroll="1"
    :mouseDrag="true"
  >
    <Slide v-for="item in reviews" :key="item.id">
      <ReviewCard
        class="slider-review-card"
        :username="item.username"
        :comment="item.comment"
        :rating="item.rating"
        :avatarUrl="item.avatarUrl"
      />
    </Slide>
  </Carousel>
</template>

<script lang="ts" setup>
import 'vue3-carousel/dist/carousel.css';
import { useWindowSize } from '@vueuse/core';
import { computed } from 'vue';
import { Carousel, Slide } from 'vue3-carousel';

import ReviewCard from '@/components/landing/ReviewCard.vue';
import type { Review } from '@/data/landing/reviews.js';

const { width } = useWindowSize();

const itemsToShow = computed(() => {
  return width.value / 408;
});

defineProps<{
  reviews: Review[];
}>();
</script>

<style lang="postcss">
.slider-review-card {
  width: 380px;
  margin: 0 12px;

  @media screen and (max-width: 565.98px) {
    width: calc(100vw - 24px * 2);
  }
}
</style>
