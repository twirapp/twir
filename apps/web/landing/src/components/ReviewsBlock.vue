<template>
  <div class="bg-black-10">
    <div class="container">
      <h2 class="text-5xl font-semibold max-w-xl leading-[125%]">
        Reviews from streamers and other viewers
      </h2>
      <ul class="inline-grid grid-flow-col gap-x-6">
        <li v-for="item in reviews" :key="item.id">
          <ReviewCard :username="item.username" :comment="item.comment" :rating="item.rating" :avatarUrl="item.avatarUrl" />
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts" setup>
import ReviewCard from '@/components/ReviewCard.vue';

const reviews: {
  id: number;
  username: string;
  comment: string;
  avatarUrl: string;
  rating: number;
}[] = [
  {
    id: 1,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: 'src/assets/avatar.png',
    rating: 4,
  },
  {
    id: 2,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: 'src/assets/avatar.png',
    rating: 4,
  },
  {
    id: 3,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: 'src/assets/avatar.png',
    rating: 2,
  },
  {
    id: 4,
    username: 'random_user',
    comment:
      'Lorem luctus tincidunt elementum dolor. Id morbi tortor mauris, viverra eu quis et id egestas.',
    avatarUrl: 'src/assets/avatar.png',
    rating: 1,
  },
];
</script>
