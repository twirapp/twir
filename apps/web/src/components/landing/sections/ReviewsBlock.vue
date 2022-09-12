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
          <SwiperSlide v-for="item in reviews" :key="item.id" style="width: 380px">
            <ReviewCard
              :username="item.username"
              :comment="item.comment"
              :rating="item.rating"
              :avatarUrl="item.avatarUrl"
            />
          </SwiperSlide>
        </Swiper>
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
import { Autoplay, type Swiper as ISwiper } from 'swiper';
import { Swiper, SwiperSlide } from 'swiper/vue';
import { ref } from 'vue';

import 'swiper/css';

import PurpleBlob from '@/assets/blob-purple.png';
import ClientOnly from '@/components/ClientOnly.vue';
import ReviewCard from '@/components/landing/ReviewCard.vue';
import useTranslation from '@/hooks/useTranslation.js';
import type { Review } from '@/types/review';

defineProps<{reviews: Review[]}>();

const slider = ref<ISwiper | null>(null);

const t = useTranslation<'landing'>();

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
