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
  <div class="masonry sm:masonry-sm md:masonry-md lg:masonry-lg m-3">
    <Spotify />

    <div class="grid grid-cols-1 mb-[0.5rem] break-inside">
      <div><Vk /></div>
      <div>        <Lastfm />    </div>
    </div>

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
