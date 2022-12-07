<script lang="ts" setup>
import { mdiLinkBoxVariant, mdiCapsLock, mdiClipboard, mdiCalculatorVariant, mdiEmoticon, mdiWrapDisabled } from '@mdi/js';
import { useDisplay } from 'vuetify';

import Card from '@/components/moderation/card.vue';

const { smAndDown } = useDisplay();
</script>

<template>
  <div>
    <v-row justify="center">
      <v-col :cols="smAndDown ? 12 : 4">
        <Card type="links" :icon="mdiLinkBoxVariant" icon-color="blue" />
      </v-col>
      <v-col :cols="smAndDown ? 12 : 4">
        <Card type="caps" :icon="mdiCapsLock" icon-color="orange" />
      </v-col>
    </v-row>
    <v-row justify="center">
      <v-col :cols="smAndDown ? 12 : 4">
        <Card type="symbols" :icon="mdiCalculatorVariant" icon-color="green" />
      </v-col>

      <v-col :cols="smAndDown ? 12 : 4">
        <Card type="emotes" :icon="mdiEmoticon" icon-color="yellow" />
      </v-col>
    </v-row>
    <v-row justify="center">
      <v-col :cols="smAndDown ? 12 : 4">
        <Card type="blacklists" :icon="mdiClipboard" icon-color="red" />
      </v-col>
      <v-col :cols="smAndDown ? 12 : 4">
        <Card type="longMessage" :icon="mdiWrapDisabled" icon-color="cyan" />
      </v-col>
    </v-row>
  </div>
</template>
