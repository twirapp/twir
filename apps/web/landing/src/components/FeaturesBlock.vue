<template>
  <div class="py-[100px] bg-black-10">
    <div class="container max-w-[1040px]">
      <div class="flex w-full justify-between mb-14">
        <h2 class="font-semibold text-5xl leading-[130%]">
          Bot features
        </h2>
        <a class="inline-flex items-center">
          <span class="text-right text-xl text-gray-70 pr-[14px]">
            Features in<br />
            development
          </span>
          <TswArrowIcon arrowType="in-circle" :strokeWidth="1.5" size="52px" direction="top-right" class="stroke-gray-35" />
        </a>
      </div>
      <ul class="grid gap-6 grid-cols-2">
        <li v-for="item in featuresInfo" :key="item.id">
          <FeatureCard :title="item.title" :description="item.description" :actionText="item.actionText" :actionHref="item.actionHref" />
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { TswArrowIcon } from '@tsuwari/ui-components';

import FeatureCard from '@/components/FeatureCard.vue';

const featuresInfo: {
  title: string;
  description: string;
  actionText: string;
  actionHref: string;
  id: number;
}[] = [
  {
    id: 1,
    title: 'Commands',
    description:
      'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
    actionText: 'Try feature',
    actionHref: '#',
  },
  {
    id: 2,
    title: 'Moderation',
    description:
      'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
    actionText: 'Try feature',
    actionHref: '#',
  },
  {
    id: 3,
    title: 'Timers',
    description:
      'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
    actionText: 'Try feature',
    actionHref: '#',
  },
  {
    id: 4,
    title: 'Greatings',
    description:
      'Facilisi eget laoreet quam fringilla pulvinar diam. Risus massa ut pellentesque mi. Facilisi lobortis erat nibh diam nunc nunc. Sit natoque risus, ut malesuada',
    actionText: 'Try feature',
    actionHref: '#',
  },
];
</script>
