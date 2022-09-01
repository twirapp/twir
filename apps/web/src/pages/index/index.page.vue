<template>
  <div class="relative min-h-screen">
    <Header :menuItems="tm('navMenu')" />
    <main>
      <FirstScreen />
      <StatsLine :stats="stats" />
      <FeaturesBlock :id="navMenuHrefs[NavMenuTabs.features]" :features="features" />
      <IntegrationsBlock />
      <ReviewsBlock :id="navMenuHrefs[NavMenuTabs.reviews]" :reviews="reviews" />
      <TeamBlock :id="navMenuHrefs[NavMenuTabs.team]" :teamMembers="teamMembers" />
      <PricingBlock
        :id="navMenuHrefs[NavMenuTabs.pricing]"
        :planColorThemes="planColorThemes"
        :pricePlans="pricePlans"
      />
    </main>
  </div>
  <Footer :menuItems="tm('navMenu')" :socials="socials" />
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import Footer from '@/components/landing/layout/Footer.vue';
import Header from '@/components/landing/layout/Header.vue';
import FeaturesBlock from '@/components/landing/sections/FeaturesBlock.vue';
import FirstScreen from '@/components/landing/sections/FirstScreen.vue';
import IntegrationsBlock from '@/components/landing/sections/IntegrationsBlock.vue';
import PricingBlock from '@/components/landing/sections/PricingBlock.vue';
import ReviewsBlock from '@/components/landing/sections/ReviewsBlock.vue';
import StatsLine from '@/components/landing/sections/StatsLine.vue';
import TeamBlock from '@/components/landing/sections/TeamBlock.vue';
import {
  reviews,
  stats,
  features,
  planColorThemes,
  pricePlans,
  teamMembers,
  navMenuHrefs,
} from '@/data';
import { socials } from '@/data';
import { NavMenuTabs } from '@/types/navMenu';

const { tm } = useI18n();
</script>
