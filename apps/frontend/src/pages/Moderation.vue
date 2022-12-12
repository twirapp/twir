<script lang="ts" setup>
import { mdiLinkBoxVariant, mdiCapsLock, mdiClipboard, mdiCalculatorVariant, mdiEmoticon, mdiWrapDisabled } from '@mdi/js';
import chunk from 'lodash.chunk';
import { useDisplay } from 'vuetify';

import Card from '@/components/moderation/card.vue';

const { smAndDown } = useDisplay();

const types = chunk([
  { key: 'links', icon: mdiLinkBoxVariant, iconColor: 'blue' }, 
  { key: 'caps', icon: mdiCapsLock,  iconColor: 'orange' },
  { key: 'symbols', icon: mdiCalculatorVariant, iconColor: 'green' }, 
  { key: 'emotes', icon: mdiEmoticon, iconColor: 'yellow' },
  { key: 'blacklists', icon: mdiClipboard, iconColor: 'red' }, 
  { key: 'longMessage', icon: mdiWrapDisabled, iconColor: 'cyan' },
], 2);
</script>

<template>
  <div>
    <v-row
      v-for="(items, typesIndex) of types"
      :key="typesIndex"
      justify="center"
    >
      <v-col v-for="(item, itemsIndex) of items" :key="itemsIndex" :cols="smAndDown ? 12 : 4">
        <Card :type="(item.key as any)" :icon="item.icon" :icon-color="item.iconColor" />
      </v-col>
    </v-row>
  </div>
</template>
