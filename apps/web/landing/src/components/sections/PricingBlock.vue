<template>
  <section class="bg-black-10 relative overflow-hidden">
    <div class="container relative z-0 py-24">
      <div class="flex flex-col items-center">
        <h2 class="text-5xl font-semibold leading-[130%] max-w-xl text-center mb-16">
          We’ve got a plan that’s perfect for you
        </h2>
        <ul class="max-w-[880px] w-full mx-auto grid grid-flow-col gap-x-[30px]">
          <li v-for="plan in pricePlans" :key="plan.name" class="flex w-full">
            <PricingPlan :plan="plan" :colorTheme="planColorThemes[plan.id]" />
          </li>
        </ul>
      </div>
      <img src="@/assets/Waves.svg" class="absolute -z-[1] -bottom-[270px] -right-[300px]" />
      <img src="@/assets/BlobCyan.svg" class="absolute -z-[1] -top-[600px] -left-[280px]" />
    </div>
  </section>
</template>

<script lang="ts" setup>
import PricingPlan from '@/components/PricingPlan.vue';
import type { PlanColorThemes } from '@/types/pricingPlan';
import type { PricePlan } from '@/types/pricingPlan';

defineProps<{pricePlans: PricePlan[], planColorThemes: PlanColorThemes}>();
</script>
