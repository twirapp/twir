<template>
  <div class="bg-black-10 relative overflow-hidden z-0">
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
  </div>
</template>

<script lang="ts" setup>
import { Autoplay } from 'swiper';
import { Swiper, SwiperSlide } from 'swiper/vue';

import TextAvatarUrl from '@/assets/avatar.png';
import ReviewCard from '@/components/ReviewCard.vue';

const modules = [Autoplay];

import 'swiper/css';

const reviews: {
  id: number;
  username: string;
  comment: string;
  avatarUrl: string;
  rating: number;
}[] = [
  {
    id: 1,
    username: 'random_usergsdagdsagsadgsda',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 4,
  },
  {
    id: 2,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 4,
  },
  {
    id: 3,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 2,
  },
  {
    id: 4,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 1,
  },
  {
    id: 5,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: TextAvatarUrl,
    rating: 4,
  },
  {
    id: 6,
    username: 'random_user',
    comment:
      'Praesent dolor quis aliquam nulla id in orci. Mi sit pulvinar nunc blandit egestas cras. Sed habitant amet ultrices vitae. At volutpat enim vel quam dignissim ut justo.',
    avatarUrl: TextAvatarUrl,
    rating: 4,
  },
];
</script>
