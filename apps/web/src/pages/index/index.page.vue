<template>
  <FirstScreen />
  <StatsLine :stats="stats" />
  <FeaturesBlock
    :id="navMenuHrefs[NavMenuTabs.features]"
    :features="tm('sections.features.content')"
  />
  <IntegrationsBlock />
  <ReviewsBlock :id="navMenuHrefs[NavMenuTabs.reviews]" :reviews="reviews" />
  <TeamBlock :id="navMenuHrefs[NavMenuTabs.team]" :teamMembers="teamMembers" />
  <PricingBlock
    :id="navMenuHrefs[NavMenuTabs.pricing]"
    :pricePlans="tm('sections.pricing.plans')"
  />
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import FeaturesBlock from '@/components/landing/sections/FeaturesBlock.vue';
import FirstScreen from '@/components/landing/sections/FirstScreen.vue';
import IntegrationsBlock from '@/components/landing/sections/IntegrationsBlock.vue';
import PricingBlock from '@/components/landing/sections/PricingBlock.vue';
import ReviewsBlock from '@/components/landing/sections/ReviewsBlock.vue';
import StatsLine from '@/components/landing/sections/StatsLine.vue';
import TeamBlock from '@/components/landing/sections/TeamBlock.vue';
import { reviews, stats, teamMembers, navMenuHrefs } from '@/data';
import { NavMenuTabs } from '@/types/navMenu';

const { tm } = useI18n();
</script>
