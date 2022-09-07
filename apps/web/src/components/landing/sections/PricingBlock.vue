<template>
  <section class="bg-black-10 relative overflow-hidden">
    <div class="container relative z-0 py-24">
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
      <ClientOnly>
        <BlurryBlob class="w-[324px] h-[400px] rotate-90 -top-[250px] left-[60px]" color="cyan" />
      </ClientOnly>
    </div>
  </section>
</template>

<script lang="ts" setup>
import BlurryBlob from '@/components/BlurryBlob.vue';
import ClientOnly from '@/components/ClientOnly.vue';
import PricingPlan from '@/components/landing/PricingPlan.vue';
import { planColorThemes } from '@/data/pricingPlans.js';
import type { PricePlansLocale } from '@/types/pricingPlan';
import { useTranslation } from '@/utils/locales.js';

defineProps<{
  pricePlans: PricePlansLocale;
}>();

const t = useTranslation<'landing'>();
</script>
