<template>
  <section class="bg-black-10 relative overflow-hidden min-md:px-10 min-sm:px-8">
    <div class="container relative z-0 py-24 max-sm:pb-0 max-sm:pt-18 max-w-[1200px]">
      <div class="flex flex-col items-center">
        <h2
          class="
            min-xl:text-[48px] min-md:text-[44px] min-sm:text-[40px]
            text-[36px]
            font-semibold
            leading-[130%]
            min-xl:max-w-[41rem] min-md:max-w-xl min-sm:max-w-md
            text-center
            mb-16
            tracking-tight
            max-sm:px-6
          "
        >
          {{ t('sections.pricing.title') }}
        </h2>
        <ul
          class="
            max-w-[880px]
            w-full
            mx-auto
            grid
            min-md:grid-flow-col
            gap-x-[30px]
            min-sm:gap-y-7 min-md:grid-cols-2
          "
        >
          <li v-for="(plan, planId) in pricePlans" :key="planId" class="flex w-full">
            <PricingPlan
              :plan="plan"
              :colorTheme="planColorThemes[planId]"
              :planGeneral="planFeaturesGeneral[planId]"
            />
          </li>
        </ul>
      </div>
      <div
        class="
          bg-no-repeat
          absolute
          -z-[1]
          -bottom-[270px]
          -right-[240px]
          bg-contain
          h-[780px]
          w-[774px]
          max-sm:hidden
        "
        :style="{
          backgroundImage: cssURL(WavesSvg),
        }"
      ></div>
      <ClientOnly>
        <TswImage
          :src="CyanBlob"
          renderType="bg-image"
          class="blurry-blob -top-[580px] -left-[540px]"
          :height="1021"
          :width="1102"
          :lazy="true"
        />
      </ClientOnly>
    </div>
  </section>
</template>

<script lang="ts" setup>
import { cssURL, TswImage } from '@twir/ui-components';
import { computed } from 'vue';

import CyanBlob from '@/assets/blob-cyan.png';
import WavesSvg from '@/assets/Waves.svg';
import ClientOnly from '@/components/ClientOnly.vue';
import PricingPlan from '@/components/landing/PricingPlan.vue';
import { planColorThemes, planFeaturesGeneral } from '@/data/landing/pricingPlans.js';
import { useTranslation } from '@/services/locale';

const { t, tm } = useTranslation<'landing'>();

const pricePlans = computed(() => tm('sections.pricing.plans'));
</script>
