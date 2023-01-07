<template>
  <HeroSection :id="LandingSection[LandingSection.hero]" />
  <StatInfoSection />
  <FeaturesSection :id="LandingSection[LandingSection.features]" />
  <IntegrationsSection />
  <ReviewsSection :id="LandingSection[LandingSection.reviews]" />
  <TeamSection :id="LandingSection[LandingSection.team]" />
  <PricingSection :id="LandingSection[LandingSection.pricing]" />
</template>

<script lang="ts" setup>
import FeaturesSection from '@/components/landing/sections/FeaturesSection.vue';
import HeroSection from '@/components/landing/sections/HeroSection.vue';
import IntegrationsSection from '@/components/landing/sections/IntegrationsSection.vue';
import PricingSection from '@/components/landing/sections/PricingSection.vue';
import ReviewsSection from '@/components/landing/sections/ReviewsSection.vue';
import StatInfoSection from '@/components/landing/sections/StatInfoSection.vue';
import TeamSection from '@/components/landing/sections/TeamSection.vue';
import { LandingSection } from '@/data/landing/sections.js';
</script>
