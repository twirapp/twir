<script lang="ts" setup>
import DonationAlerts from '@/components/integrations/donationalerts.vue';
import Faceit from '@/components/integrations/faceit.vue';
import Lastfm from '@/components/integrations/lastfm.vue';
import Spotify from '@/components/integrations/spotify.vue';
import Streamlabs from '@/components/integrations/streamlabs.vue';
import VK from '@/components/integrations/vk.vue';
</script>

<template>
  <div>
    <div class="d-flex justify-center flex-column flex-sm-row flex-md-row flex-lg-row flex-xl-row">
      <Spotify />
      <VK />
      <Lastfm />
    </div>
    <div class="d-flex justify-center flex-column flex-sm-row flex-md-row flex-lg-row flex-xl-row">
      <Faceit />
      <Streamlabs />
      <DonationAlerts />
    </div>
  </div>
</template>