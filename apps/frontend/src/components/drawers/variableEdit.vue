<script lang="ts" setup>
import { mdiPlus, mdiClose } from '@mdi/js';
import { type ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { ref } from 'vue';

import confirmDeletion from '@/components/confirmDeletion.vue';

const props = defineProps<{
  variable: ChannelCustomvar
}>();

const emits = defineEmits<{
  (event: 'cancel'): () => void
}>();

const variable = ref(props.variable);

function onDelete() {
  console.log(variable);
  emits('cancel');
}
</script>

<template>
  <div>
    <v-list-item>
      <div class="d-flex justify-space-between">
        <h1>Edit timer</h1>
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
          <v-text-field 
            v-model="variable.name" 
            label="Name" 
            :rules="[
              v => !!v || 'Field is required'
            ]"
          />
        </div>
      </v-form>
    </v-list-item>
  </div>
</template>