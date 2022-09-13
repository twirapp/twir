<template>
  <section class="bg-black-10 relative overflow-hidden z-0">
    <div class="container relative max-w-[980px] mt-20 mb-16">
      <img src="@/assets/CurveArrow.svg" class="absolute top-[30px] -left-[90px] -z-[1]" />
      <div
        :style="{
          backgroundImage: `url('${PurpleBlob}')`,
        }"
        class="blurry-blob -top-[550px] -right-[500px] w-[1100px] h-[1075px]"
      ></div>
      <img src="@/assets/MessageCircle.svg" class="absolute -z-[1] -bottom-3 right-28" />
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
import ClientOnly from '@/components/ClientOnly.vue';
import ReviewCard from '@/components/landing/ReviewCard.vue';
import ReviewsSlider from '@/components/landing/ReviewsSlider.vue';
import { reviews } from '@/data/landing/reviews.js';
import useTranslation from '@/hooks/useTranslation.js';

const t = useTranslation<'landing'>();
</script>
