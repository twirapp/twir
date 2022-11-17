<template>
  <div :class="`inline-flex flex-col min-sm:rounded-[10px] w-full ${colorTheme}-pricing-plan`">
    <div class="plan-info p-6 grid gap-y-3 justify-items-start border-b">
      <h3 class="text-[21px] leading-[130%]">
        {{ plan.name }}
      </h3>
      <span class="inline-grid grid-flow-col items-baseline gap-x-2">
        <span class="text-[44px] font-medium">${{ plan.price }}</span>
        <span class="price-per">{{ t('sections.pricing.perMonth') }}</span>
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
          transition-colors
        "
      >
        {{ buttonText }}
      </a>
    </div>
    <div class="p-6 max-sm:pb-9">
      <span class="uppercase tracking-wider plan-features-title mb-5 block">
        {{ t('sections.pricing.features') }}
      </span>
      <ul class="grid gap-y-3">
        <li
          v-for="(feature, featureId) in plan.features"
          :key="featureId"
          class="inline-grid grid-flow-col justify-start"
        >
          <TswIcon
            :name="featureTypeIcons[feature.status]"
            :height="22"
            :width="22"
            :class="`feature-icon mr-3 ${
              colorTheme === 'purple'
                ? 'stroke-white-100'
                : feature.status === 'accessible'
                ? 'stroke-purple-80'
                : 'stroke-gray-70'
            }`"
          />
          <span class="w-full">{{ feature.name }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { TswIcon } from '@tsuwari/ui-components';
import { computed } from 'vue';

import { featureTypeIcons } from '@/data/landing/pricingPlans.js';
import type { PlanColorTheme, PricePlanLocale } from '@/data/landing/pricingPlans.js';
import { useTranslation } from '@/services/locale';

const props =
  defineProps<{
    plan: PricePlanLocale;
    colorTheme: PlanColorTheme;
  }>();

const { t } = useTranslation<'landing'>();

const buttonText = computed(() =>
  props.plan.price > 0 ? t('buttons.buyPlan') : t('buttons.getStarted'),
);
</script>

<style lang="postcss">
.gray-pricing-plan {
  @apply bg-black-15;

  .price-per {
    @apply text-gray-70;
  }

  .action-button {
    @apply bg-purple-60 hover:bg-purple-50;
  }

  .plan-info {
    @apply border-b-black-25;
  }

  .plan-features-title {
    @apply text-gray-70;
  }
}

.purple-pricing-plan {
  background: linear-gradient(180deg, #513ada 0%, #522cbd 100%);

  .price-per {
    @apply text-white-100 text-opacity-70;
  }

  .action-button {
    @apply bg-purple-95 text-purple-55 font-medium hover:bg-opacity-0 hover:text-white-100 border hover:border-white-95;
  }

  .plan-info {
    @apply border-b-white-100 border-opacity-30;
  }

  .plan-features-title {
    @apply text-white-100 text-opacity-70;
  }
}
</style>
