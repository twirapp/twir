<script lang="ts" setup>
import { mdiPlus, mdiClose } from '@mdi/js';
import { ChannelModerationSetting } from '@tsuwari/typeorm/entities/ChannelModerationSetting';
import { ref } from 'vue';

import confirmDeletion from '@/components/confirmDeletion.vue';

const props = defineProps<{
  settings: ChannelModerationSetting
}>();

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
        <h1>{{ settings.type.charAt(0).toUpperCase() + settings.type.substring(1) }}</h1>
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
        <div
          class="d-flex flex-column"
        >
          qwe
        </div>
      </v-form>
    </v-list-item>
  </div>
</template>