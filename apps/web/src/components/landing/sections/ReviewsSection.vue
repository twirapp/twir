<template>
  <section class="bg-black-10 relative overflow-hidden z-0">
    <div class="px-8">
      <div class="container relative max-w-[980px] mt-20 mb-16">
        <div
          class="
            absolute
            top-[30px]
            min-xl:-left-[90px]
            -z-[1]
            w-[48px]
            h-[105px]
            min-lg:left-[160px] min-md:left-[100px] min-sm:left-[20px]
            -left-[22px]
            max-sm:top-[12px] max-sm:scale-110
          "
          :style="{
            backgroundImage: cssURL(CurveArrow),
          }"
        ></div>
        <div
          :style="{
            backgroundImage: cssURL(PurpleBlob),
          }"
          class="
            blurry-blob
            -top-[550px]
            -right-[500px]
            w-[1100px]
            h-[1075px]
            bg-center
            max-lg:-right-[50%] max-lg:-left-[50%] max-lg:mx-auto max-lg:w-[760px]
            bg-contain
            max-lg:-top-[200px]
          "
        ></div>
        <div
          class="absolute -z-[1] -bottom-3 right-[11%] max-lg:hidden w-[91px] h-[88px]"
          :style="{
            backgroundImage: cssURL(MessageCircle),
          }"
        ></div>
        <h2
          class="
            text-5xl
            max-w-xl
            leading-[125%]
            tracking-tight
            min-xl:text-[48px] min-sm:text-[44px]
            text-[36px]
            font-semibold
            max-xl:text-center max-xl:max-w-md max-xl:mx-auto
          "
        >
          {{ t('sections.reviews.title') }}
        </h2>
      </div>
    </div>

    <ClientOnly>
      <template #default>
        <ReviewsSlider :reviews="reviews" />
      </template>
      <template #server>
        <div
          class="
            inline-grid
            grid-flow-col
            w-screen
            overflow-hidden
            mb-20
            gap-x-6
            max-md:px-6
            min-md:justify-center
          "
        >
          <div v-for="item in reviews" :key="item.id" class="slider-review-card">
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
import { cssURL } from '@/utils/css';

const t = useTranslation<'landing'>();
</script>
