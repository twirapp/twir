<template>
  <section class="bg-black-10 relative overflow-hidden z-0">
    <div class="container relative max-w-[980px] mt-20 mb-16">
      <div
        class="absolute top-[30px] -left-[90px] -z-[1] w-[48px] h-[105px]"
        :style="{
          backgroundImage: `url('${CurveArrow}')`,
        }"
      ></div>
      <div
        :style="{
          backgroundImage: `url('${PurpleBlob}')`,
        }"
        class="blurry-blob -top-[550px] -right-[500px] w-[1100px] h-[1075px]"
      ></div>
      <div
        class="absolute -z-[1] -bottom-3 right-28 w-[91px] h-[88px]"
        :style="{
          backgroundImage: `url('${MessageCircle}')`,
        }"
      ></div>
      <h2 class="text-5xl font-semibold max-w-xl leading-[125%]">
        {{ t('sections.reviews.title') }}
      </h2>
    </div>
    <ClientOnly>
      <template #default>
        <ReviewsSlider :reviews="reviews" />
      </template>
      <template #server>
        <div class="inline-grid grid-flow-col w-screen overflow-hidden mb-20 gap-x-6">
          <div v-for="item in reviews" :key="item.id" style="width: 380px">
            <ReviewCard
              :username="item.username"
              :comment="item.comment"
              :rating="item.rating"
              :avatarUrl="item.avatarUrl"
            />
          </div>
        </div>
      </template>
    </ClientOnly>
  </section>
</template>

<script lang="ts" setup>
import PurpleBlob from '@/assets/blob-purple.png';
import CurveArrow from '@/assets/CurveArrow.svg';
import MessageCircle from '@/assets/MessageCircle.svg';
import ClientOnly from '@/components/ClientOnly.vue';
import ReviewCard from '@/components/landing/ReviewCard.vue';
import ReviewsSlider from '@/components/landing/ReviewsSlider.vue';
import { reviews } from '@/data/landing/reviews.js';
import useTranslation from '@/hooks/useTranslation.js';

const t = useTranslation<'landing'>();
</script>
