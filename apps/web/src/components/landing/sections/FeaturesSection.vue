<template>
  <section class="py-24 max-md:py-18 bg-black-10 min-md:px-10">
    <div class="container max-w-[1200px]">
      <div class="flex w-full justify-between mb-14 max-md:mb-7 max-md:flex-col">
        <h2
          class="
            font-semibold
            text-5xl
            max-md:text-[42px]
            leading-[130%]
            tracking-tight
            min-md:px-0 min-sm:px-8
            max-sm:px-5 max-md:mb-8
          "
        >
          {{ t('sections.features.title') }}
        </h2>
        <a
          class="
            inline-flex
            items-center
            max-md:w-full max-md:bg-purple-60 max-md:py-3
            min-md:px-0 min-sm:px-8
            px-6
            max-md:justify-end
          "
          href="#"
        >
          <span
            class="
              text-right text-xl
              max-md:text-[18px] max-md:leading-[135%]
              text-gray-70
              pr-[14px]
              max-w-[140px] max-md:text-white-100
            "
          >
            {{ t('sections.features.featuresInDev') }}
          </span>
          <TswArrowIcon
            arrowName="ArrowInCircle"
            :strokeWidth="1.5"
            size="52px"
            direction="top-right"
            class="stroke-gray-35 max-md:stroke-white-100"
          />
        </a>
      </div>
      <ul
        class="
          grid
          min-md:gap-6
          max-md:gap-x-6
          grid-cols-2
          max-sm:grid-cols-1
          min-md:px-0 min-sm:px-8
          max-sm:px-5
        "
      >
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
