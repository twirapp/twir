<script lang="ts" setup>
import { useDisplay } from 'vuetify';

import DonationAlerts from '@/components/integrations/donationalerts.vue';
import Faceit from '@/components/integrations/faceit.vue';
import Lastfm from '@/components/integrations/lastfm.vue';
import Spotify from '@/components/integrations/spotify.vue';
import Streamlabs from '@/components/integrations/streamlabs.vue';
import VK from '@/components/integrations/vk.vue';

const { smAndDown } = useDisplay();
</script>

<template>
  <div>
    <v-row justify="center">
      <v-col :cols="smAndDown ? 12 : 3">
        <Spotify />
      </v-col>
      <v-col :cols="smAndDown ? 12 : 3">
        <VK />
      </v-col>
      <v-col :cols="smAndDown ? 12 : 3">
        <Lastfm />
      </v-col>
    </v-row>
    <v-row justify="center">
      <v-col :cols="smAndDown ? 12 : 3">
        <Faceit />
      </v-col>
      <v-col :cols="smAndDown ? 12 : 3">
        <Streamlabs />
      </v-col>
      <v-col :cols="smAndDown ? 12 : 3">
        <DonationAlerts />
      </v-col>
    </v-row>
  </div>
</template>