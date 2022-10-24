<template>
  <Flicking
    :options="{
      align: 'center',
      renderOnlyVisible: true,
      moveType: 'snap',
      circular: true,
      circularFallback: 'bound',
    }"
    class="mb-20 cursor-grab"
    :plugins="plugins"
  >
    <ReviewCard
      v-for="item in reviews"
      :key="item.id"
      class="slider-review-card"
      :username="item.username"
      :comment="item.comment"
      :rating="item.rating"
      :avatarUrl="item.avatarUrl"
    />
  </Flicking>
</template>

<script lang="ts" setup>
import { AutoPlay } from '@egjs/flicking-plugins';
import Flicking from '@egjs/vue3-flicking';

import ReviewCard from '@/components/landing/ReviewCard.vue';
import type { Review } from '@/data/landing/reviews.js';

defineProps<{
  reviews: Review[];
}>();

const plugins = [new AutoPlay({ stopOnHover: true })];
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

<style>
@import '@egjs/vue3-flicking/dist/flicking.css';
</style>
