<script lang="ts" setup>
import { type SettingsType } from '@tsuwari/typeorm/entities/ChannelModerationSetting';

defineProps<{
  type: keyof typeof SettingsType,
  icon: string,
  iconColor: string,
}>();
</script>

<template>
  <v-card class="ma-2 d-flex flex-column">
    <v-card-title>
      <div class="d-flex justify-space-between">
        <div>
          <v-icon :icon="icon" :color="iconColor"></v-icon>
          {{ type.charAt(0).toUpperCase() + type.substring(1) }}
        </div>
        <div style="height: 10px;">
          <v-switch style="margin-top:-10px" color="indigo"></v-switch>
        </div>
      </div>
    </v-card-title>
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
  </v-card>
</template>