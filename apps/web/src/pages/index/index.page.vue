<template>
  <HeroSection />
  <StatInfoSection />
  <FeaturesSection :id="navMenuHrefs[NavMenuTabs.features]" />
  <IntegrationsSection />
  <ReviewsSection :id="navMenuHrefs[NavMenuTabs.reviews]" />
  <TeamSection :id="navMenuHrefs[NavMenuTabs.team]" />
  <PricingSection :id="navMenuHrefs[NavMenuTabs.pricing]" />
</template>

<script lang="ts" setup>
import FeaturesSection from '@/components/landing/sections/FeaturesSection.vue';
import HeroSection from '@/components/landing/sections/HeroSection.vue';
import IntegrationsSection from '@/components/landing/sections/IntegrationsSection.vue';
import PricingSection from '@/components/landing/sections/PricingSection.vue';
import ReviewsSection from '@/components/landing/sections/ReviewsSection.vue';
import StatInfoSection from '@/components/landing/sections/StatInfoSection.vue';
import TeamSection from '@/components/landing/sections/TeamSection.vue';
import { navMenuHrefs, NavMenuTabs } from '@/data/landing/navMenu.js';
</script>
