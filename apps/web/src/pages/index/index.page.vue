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
import FeaturesBlock from '@/components/sections/FeaturesBlock.vue';
import FirstScreen from '@/components/sections/FirstScreen.vue';
import IntegrationsBlock from '@/components/sections/IntegrationsBlock.vue';
import PricingBlock from '@/components/sections/PricingBlock.vue';
import ReviewsBlock from '@/components/sections/ReviewsBlock.vue';
import StatsLine from '@/components/sections/StatsLine.vue';
import TeamBlock from '@/components/sections/TeamBlock.vue';
import * as data from '@/data';
</script>
