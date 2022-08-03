<script lang="ts" setup>
import Setup from '../components/landing/Setup.vue';
import Team from '../components/landing/Team.vue';

import Features from '@/components/landing/Features.vue';
import Footer from '@/components/landing/Footer.vue';
import Statistics from '@/components/landing/Statistics.vue';
import Streamers from '@/components/landing/Streamers.vue';
import UseIt from '@/components/landing/UseIt.vue';
import Welcome from '@/components/landing/Welcome.vue';
import MainNavBar from '@/components/layout/MainNavBar.vue';

document.title = 'Tsuwari - Main page';
</script>

<template>
  <div style="font-family: 'Inter', sans-serif">
    <MainNavBar />

    <Welcome />

    <Statistics />
    <Features />
    <Setup />
    <Streamers />
    <UseIt />

    <Team />

    <Footer />
  </div>
</template>
