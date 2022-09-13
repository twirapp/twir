<template>
  <section class="py-[100px] bg-black-10">
    <div class="container max-w-[1040px]">
      <div class="flex w-full justify-between mb-14">
        <h2 class="font-semibold text-5xl leading-[130%]">
          {{ t('sections.features.title') }}
        </h2>
        <a class="inline-flex items-center">
          <span class="text-right text-xl text-gray-70 pr-[14px] max-w-[140px]">
            {{ t('sections.features.featuresInDev') }}
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
        <li v-for="(feature, i) in features" :key="i">
          <FeatureCard :title="feature.name" :description="feature.description" :actionHref="'#'" />
        </li>
      </ul>
    </div>
  </section>
</template>

<script lang="ts" setup>
import { TswArrowIcon } from '@tsuwari/ui-components';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import FeatureCard from '@/components/landing/FeatureCard.vue';
import type { BotFeature } from '@/data/landing/botFeatures';
import useTranslation from '@/hooks/useTranslation.js';

const { tm } = useI18n();

const features = computed(() => tm('sections.features.content') as BotFeature[]);

const t = useTranslation<'landing'>();
</script>
