<script lang="ts" setup>
import { 
  mdiLinkBoxVariant, 
  mdiCapsLock, 
  mdiClipboard, 
  mdiCalculatorVariant, 
  mdiEmoticon, 
  mdiWrapDisabled, 
  mdiPencil, 
} from '@mdi/js';
import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import chunk from 'lodash.chunk';
import { ref } from 'vue';
import { useDisplay } from 'vuetify';

import ModerationEdit from '@/components/drawers/moderationEdit.vue';
import { moderationSettings } from '@/stores/moderationSettings';

const { smAndDown } = useDisplay();
const editableSettings = ref<ChannelModerationSetting | undefined>();
const isEditDrawerOpened = ref(false);

const types = chunk([
  { key: 'links', icon: mdiLinkBoxVariant, iconColor: 'blue' }, 
  { key: 'caps', icon: mdiCapsLock,  iconColor: 'orange' },
  { key: 'symbols', icon: mdiCalculatorVariant, iconColor: 'green' }, 
  { key: 'emotes', icon: mdiEmoticon, iconColor: 'yellow' },
  { key: 'blacklists', icon: mdiClipboard, iconColor: 'red' }, 
  { key: 'longMessage', icon: mdiWrapDisabled, iconColor: 'cyan' },
], 2);

const descriptions = {
  'links': `Remove messages containing any links you haven't whitelisted.`,
  'caps': `Remove messages containing excessive amounts of capital letters.`,
  'symbols': `Remove messages containing disruptive or excessive use of symbols.`,
  'longMessage': `Remove lengthy messages.`,
  'emotes': 'Remove messages containing an excessive amount of emotes.',
  'blacklists': 'Remove blacklisted words.',
} as { [x: string]: string };

function cancelEdit() {
  isEditDrawerOpened.value = false;
  editableSettings.value = undefined;
}

function setEditSettings(key: string) {
  const settings = moderationSettings.value.find(s => s.type === key);
  if (!settings) return;

  isEditDrawerOpened.value = true;
  editableSettings.value = JSON.parse(JSON.stringify(settings));
}
</script>

<template>
  <div>
    <v-row
      v-for="(items, typesIndex) of types"
      :key="typesIndex"
      justify="center"
    >
      <v-col v-for="(item, itemsIndex) of items" :key="itemsIndex" :cols="smAndDown ? 12 : 4">
        <v-card class="ma-2 d-flex flex-column">
          <v-card-title>
            <div class="d-flex justify-space-between">
              <div>
                <v-icon :icon="item.icon" :color="item.iconColor"></v-icon>
                {{ item.key.charAt(0).toUpperCase() + item.key.substring(1) }}
              </div>
              <div style="height: 10px;" class="d-flex justify-space-between">
                <v-btn 
                  :icon="mdiPencil" 
                  variant="flat" 
                  size="x-small" 
                  @click="setEditSettings(item.key)" 
                />
              </div>
            </div>
          </v-card-title>
          <v-divider />
          <v-card-text>
            <v-row align="center" no-gutters>
              {{ descriptions[item.key as string] ?? '' }}
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <v-navigation-drawer
      v-if="isEditDrawerOpened"
      v-model="isEditDrawerOpened"
      temporary
      location="right"
      :class="[smAndDown ? 'w-100' : 'w-50']"
    >
      <ModerationEdit :settings="editableSettings!" @cancel="cancelEdit" />
    </v-navigation-drawer>
  </div>
</template>
