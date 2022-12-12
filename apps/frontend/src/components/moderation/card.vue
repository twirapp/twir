<script lang="ts" setup>
import { mdiPencil } from '@mdi/js';
import { type SettingsType } from '@tsuwari/typeorm/entities/ChannelModerationSetting';

defineProps<{
  type: keyof typeof SettingsType,
  icon: string,
  iconColor: string,
}>();

const descriptions = {
  'links': `Remove messages containing any links you haven't whitelisted.`,
  'caps': `Remove messages containing excessive amounts of capital letters.`,
  'symbols': `Remove messages containing disruptive or excessive use of symbols.`,
  'longMessage': `Remove lengthy messages.`,
  'emotes': 'Remove messages containing an excessive amount of emotes.',
  'blacklists': 'Remove blacklisted words.',
} as { [x: string]: string };
</script>

<template>
  <v-card class="ma-2 d-flex flex-column">
    <v-card-title>
      <div class="d-flex justify-space-between">
        <div>
          <v-icon :icon="icon" :color="iconColor"></v-icon>
          {{ type.charAt(0).toUpperCase() + type.substring(1) }}
        </div>
        <div style="height: 10px;" class="d-flex justify-space-between">
          <v-btn class="mr-2" :icon="mdiPencil" variant="tonal" size="x-small" />
          <v-switch style="margin-top:-10px" color="indigo" />
        </div>
      </div>
    </v-card-title>
    <v-divider />
    <v-card-text>
      <v-row align="center" no-gutters>
        {{ descriptions[type as string] ?? 'q' }}
      </v-row>
    </v-card-text>
  </v-card>
</template>