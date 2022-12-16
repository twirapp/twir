<script lang="ts" setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

const isBotMod = ref(false);
const { t } = useI18n();
</script>

<template>
  <v-card class="ma-2 d-flex flex-column fill-height">
    <v-card-title>Last.fm</v-card-title>
    <v-divider />
    <v-card-text>
      <v-row align="center" no-gutters>
        <v-col
          class="text-h4"
          cols="6"
        >
          Satont WorldWide
        </v-col>

        <v-col cols="6" class="text-right">
          <v-avatar 
            image="https://randomuser.me/api/portraits/women/85.jpg"
            size="100"
          />
        </v-col>
      </v-row>
    </v-card-text>
    <v-divider></v-divider>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn variant="outlined" color="primary">
        Login
      </v-btn>
    </v-card-actions>
  </v-card>
</template>