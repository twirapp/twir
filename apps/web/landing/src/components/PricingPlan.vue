<template>
  <div :class="`inline-flex flex-col rounded-[10px] w-full ${colorTheme}-plan`">
    <div class="plan-info p-6 grid gap-y-3 justify-items-start border-b">
      <h6 class="text-[21px] leading-[130%]">
        {{ plan.name }}
      </h6>
      <span class="inline-grid grid-flow-col items-baseline gap-x-2">
        <span class="text-[44px] font-medium">${{ plan.price }}</span>
        <span class="price-per">per month</span>
      </span>
      <a
        href="#"
        class="
          action-button
          inline-flex
          justify-center
          w-full
          rounded-md
          px-3
          py-[10px]
          leading-[130%]
        "
      >
        {{ buttonText }}
      </a>
    </div>
    <div class="p-6">
      <span class="uppercase tracking-wider plan-features-title mb-5 block">Features</span>
      <ul class="grid gap-y-3">
        <li v-for="(feature, index) in plan.features" :key="index" class="inline-flex">
          <TswIcon
            :name="featureTypeIcons[feature.status]"
            size="22px"
            :class="`feature-icon mr-3 ${
              colorTheme === 'purple'
                ? 'stroke-white-100'
                : feature.status === FeatureType.accessibly
                ? 'stroke-purple-80'
                : 'stroke-gray-70'
            }`"
          />
          <span>{{ feature.feature }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { TswIcon } from '@tsuwari/ui-components';
import type { IconName } from '@tsuwari/ui-icons';

import { FeatureType, PlanColorTheme } from '@/types/pricingPlan';
import type { PricePlan } from '@/types/pricingPlan';

const props = defineProps<{ plan: PricePlan; colorTheme: PlanColorTheme }>();

const featureTypeIcons: { [K in FeatureType]: IconName } = {
  [FeatureType.accessibly]: 'Check',
  [FeatureType.limited]: 'Minus',
};

const buttonText = props.plan.price > 0 ? 'Buy plan' : 'Get started';
</script>

<style lang="postcss" scoped>
.gray-plan {
  @apply bg-black-15;

  .price-per {
    @apply text-gray-70;
  }

  .action-button {
    @apply bg-purple-60;
  }

  .plan-info {
    @apply border-b-black-25;
  }

  .plan-features-title {
    @apply text-gray-70;
  }
}

.purple-plan {
  background: linear-gradient(180deg, #513ada 0%, #522cbd 100%);

  .price-per {
    @apply text-white-100 text-opacity-70;
  }

  .action-button {
    @apply bg-purple-95 text-purple-55 font-medium;
  }

  .plan-info {
    @apply border-b-white-100 border-opacity-30;
  }

  .plan-features-title {
    @apply text-white-100 text-opacity-70;
  }
}
</style>
