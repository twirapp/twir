<script lang="ts" setup>
import { mdiPlus, mdiClose } from '@mdi/js';
import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { ref } from 'vue';
import { useDisplay } from 'vuetify';

import confirmDeletion from '@/components/confirmDeletion.vue';

const props = defineProps<{
  settings: ChannelModerationSetting
}>();
const { smAndDown } = useDisplay();
const emits = defineEmits<{
  (event: 'cancel'): () => void
}>();

const settings = ref(props.settings);

function onDelete() {
  console.log(settings);
  emits('cancel');
}
</script>

<template>
  <div>
    <v-list-item>
      <div class="d-flex justify-space-between">
        <div class="d-flex flex-row">
          <h1>{{ settings.type.charAt(0).toUpperCase() + settings.type.substring(1) }}</h1>
          <v-switch 
            v-model="settings.enabled" 
            style="height: 10px; margin-top: -7px;" 
            class="ml-2"
            inset
          />
        </div>
        <div class="d-flex d-inline">
          <v-btn size="small" class="mt-1 mr-2" @click="$emit('cancel')">
            Cancel
          </v-btn>
          <confirmDeletion :cb="() => onDelete()">
            <v-btn color="red" size="small" class="mt-1 mr-2">
              Delete
            </v-btn>
          </confirmDeletion>
          <v-btn color="green" size="small" class="mt-1">
            Save
          </v-btn>
        </div>
      </div>
    </v-list-item>

    <v-divider></v-divider>

    <v-list-item>
      <v-form class="mt-2">
        <div class="d-flex flex-row">
          <v-textarea 
            v-model="settings.banMessage" 
            label="Ban message"
            auto-grow
            rows="1"
            row-height="5"
          />

          <v-textarea
            v-model.number="settings.banTime"
            class="ml-2" 
            label="Ban time"
            auto-grow
            rows="1"
            row-height="5"
          />
        </div>

        <div class="d-flex flex-row">
          <v-textarea 
            v-model="settings.warningMessage" 
            label="Warning message"
            auto-grow
            rows="1"
            row-height="5"
          />
        </div>

        <v-text-field 
          v-if="settings.type === 'blacklists'"
          v-model.number="settings.maxPercentage"
          label="Max percent of symbols in message"
        />
        <v-text-field 
          v-if="settings.type === 'longMessage'"
          v-model.number="settings.triggerLength"
          label="Max message length"
        />
        <v-text-field 
          v-if="settings.type === 'caps'"
          v-model.number="settings.maxPercentage"
          label="Max percent of caps in message"
        />
        <v-text-field 
          v-if="settings.type === 'emotes'"
          v-model.number="settings.triggerLength"
          label="Max emotes in message"
        />

        
        <v-row>
          <v-col v-if="settings.type === 'links'" :cols="smAndDown ? 12 : 4">
            <v-checkbox v-model="settings.checkClips" label="Moderate clips" density="compact" />
          </v-col>
          <v-col :cols="smAndDown ? 12 : 4">
            <v-checkbox v-model="settings.subscribers" label="Moderate subscribers" density="compact" />
          </v-col>
          <v-col :cols="smAndDown ? 12 : 4">
            <v-checkbox v-model="settings.vips" label="Moderate vips" density="compact" />
          </v-col>
        </v-row>
      </v-form>
    </v-list-item>
  </div>
</template>