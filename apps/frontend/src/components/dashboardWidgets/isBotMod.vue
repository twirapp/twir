<script lang="ts" setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const isBotMod = ref(false);
const { t } = useI18n();
</script>

<template>
  <v-card class="ma-2 w-25">
    <v-card-text>
      <v-alert density="compact" :type="isBotMod ? 'success' : 'red'">
        <span v-html="t(isBotMod ? 'pages.dashboard.widgets.status.mod' : 'pages.dashboard.widgets.status.notMod')"></span>
      </v-alert>
    </v-card-text>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn variant="outlined" color="error">
        Leave
      </v-btn>
      <v-btn variant="outlined" color="success">
        Join
      </v-btn>
    </v-card-actions>
  </v-card>
</template>