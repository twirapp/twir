<script lang="ts" setup>
import { useTitle } from '@vueuse/core';

import FaceIt from '@/components/integrations/FaceIt.vue';
import Lastfm from '@/components/integrations/Lastfm.vue';
import Spotify from '@/components/integrations/Spotify.vue';
import Vk from '@/components/integrations/Vk.vue';

const title = useTitle();
title.value = 'Tsuwari - Integrations';
</script>

<template>
  <div class="lg:masonry-lg m-1.5 mansory md:m-3 md:mansory-md sm:mansory-sm">
    <Spotify />
    <Vk />
    <Lastfm />
    <FaceIt />
  </div>
</template>

<style scoped>
@media (max-width: 426px) {
  .small-screen {
    width: 100%;
  }
}
</style>
