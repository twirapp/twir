<template>
  <section class="py-[100px] bg-black-10">
    <div class="container max-w-[1040px]">
      <div class="flex w-full justify-between mb-14">
        <h2 class="font-semibold text-5xl leading-[130%]">Bot features</h2>
        <a class="inline-flex items-center">
          <span class="text-right text-xl text-gray-70 pr-[14px]">
            Features in<br />
            development
          </span>
          <TswArrowIcon
            arrowName="ArrowInCircle"
            :strokeWidth="1.5"
            size="52px"
            direction="top-right"
            class="stroke-gray-35"
          />
        </a>
      </div>
      <ul class="grid gap-6 grid-cols-2">
        <li v-for="feature in features" :key="feature.id">
          <FeatureCard
            :title="feature.title"
            :description="feature.description"
            :actionText="feature.actionText"
            :actionHref="feature.actionHref"
          />
        </li>
      </ul>
    </div>
  </section>
</template>

<script lang="ts" setup>
import { TswArrowIcon } from '@tsuwari/ui-components';

import FeatureCard from '@/components/landing/FeatureCard.vue';
import type { BotFeature } from '@/types/botFeatures';

defineProps<{ features: BotFeature[] }>();
</script>
