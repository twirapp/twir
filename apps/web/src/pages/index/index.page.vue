<template>
  <FirstScreen />
  <StatsLine :stats="data.stats" />
  <FeaturesBlock :features="data.features" />
  <IntegrationsBlock />
  <ReviewsBlock :reviews="data.reviews" />
  <TeamBlock :teamMembers="data.teamMembers" />
  <PricingBlock :planColorThemes="data.planColorThemes" :pricePlans="data.pricePlans" />
</template>

<script lang="ts" setup>
import FeaturesBlock from '@/components/landing/sections/FeaturesBlock.vue';
import FirstScreen from '@/components/landing/sections/FirstScreen.vue';
import IntegrationsBlock from '@/components/landing/sections/IntegrationsBlock.vue';
import PricingBlock from '@/components/landing/sections/PricingBlock.vue';
import ReviewsBlock from '@/components/landing/sections/ReviewsBlock.vue';
import StatsLine from '@/components/landing/sections/StatsLine.vue';
import TeamBlock from '@/components/landing/sections/TeamBlock.vue';
import * as data from '@/data';
</script>
