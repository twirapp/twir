<script lang="ts" setup>
import { useDisplay } from 'vuetify';

import DashboardAccess from '@/components/settings/dashboardAccess.vue';

const { smAndDown } = useDisplay();
</script>

<template>
  <div>
    <v-row>
      <v-col :cols="smAndDown ? 12 : 3">
        <DashboardAccess />
      </v-col>
    </v-row>
  </div>
</template>