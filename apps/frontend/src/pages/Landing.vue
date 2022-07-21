<script lang="ts" setup>
import Setup from '../components/landing/Setup.vue';
import Statistics from '../components/landing/Statistics.vue';
import Streamers from '../components/landing/Streamers.vue';

import Welcome from '@/components/landing/Welcome.vue';
import MainNavBar from '@/components/layout/MainNavBar.vue';
</script>

<template>
  <div style="font-family: 'Inter', sans-serif">
    <MainNavBar />

    <Welcome />

    <Setup />
    <Statistics />
    <Streamers />
  </div>
</template>
