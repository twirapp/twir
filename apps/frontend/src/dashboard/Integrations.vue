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
  <div class="mansory sm:mansory-sm md:mansory-md lg:masonry-lg m-1.5 md:m-3">
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
