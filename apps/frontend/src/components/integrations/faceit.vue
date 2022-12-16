<script lang="ts" setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import FaceitIcon from '@/assets/icons/brands/faceit.svg?component';

const valid = ref(false);
const { t } = useI18n();
</script>

<template>
  <v-card class="ma-2 d-flex flex-column fill-height">
    <v-card-title>Faceit</v-card-title>
    <v-divider />
    <v-card-text>
      <v-form v-model="valid">
        <v-text-field label="Username" variant="underlined"></v-text-field>
        <v-select
          label="Select"
          :items="['CSGO']"
        ></v-select>
      </v-form>
    </v-card-text>
    <v-divider></v-divider>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-btn variant="outlined" color="primary">
        Save
      </v-btn>
    </v-card-actions>
  </v-card>
</template>