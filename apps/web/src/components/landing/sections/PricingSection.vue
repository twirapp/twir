<template>
  <section class="bg-black-10 relative overflow-hidden">
    <div class="container relative z-0 py-24 max-w-[1200px]">
      <div class="flex flex-col items-center">
        <h2 class="text-5xl font-semibold leading-[130%] max-w-[41rem] text-center mb-16">
          {{ t('sections.pricing.title') }}
        </h2>
        <ul class="max-w-[880px] w-full mx-auto grid grid-flow-col gap-x-[30px] grid-cols-2">
          <li v-for="(plan, planId) in pricePlans" :key="planId" class="flex w-full">
            <PricingPlan :plan="plan" :colorTheme="planColorThemes[planId]" />
          </li>
        </ul>
      </div>
      <img src="@/assets/Waves.svg" class="absolute -z-[1] -bottom-[270px] -right-[300px]" />
      <div
        class="blurry-blob w-[1102px] h-[1021px] -top-[580px] -left-[540px]"
        :style="{
          backgroundImage: cssURL(CyanBlob),
        }"
      ></div>
    </div>
  </section>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import CyanBlob from '@/assets/blob-cyan.png';
import PricingPlan from '@/components/landing/PricingPlan.vue';
import { planColorThemes, type PricePlansLocale } from '@/data/landing/pricingPlans.js';
import useTranslation from '@/hooks/useTranslation.js';
import { cssURL } from '@/utils/css';

const { tm } = useI18n();
const t = useTranslation<'landing'>();

const pricePlans = computed(() => tm('sections.pricing.plans') as PricePlansLocale);
</script>
