<script lang="ts" setup>
// eslint-disable-next-line import/order
import { ref } from 'vue';

import confirmDeletion from '@/components/confirmDeletion.vue';
import { Greeting } from '@/stores/greetings';

const props = defineProps<{
  greeting: Greeting
}>();
const emits = defineEmits<{
  (event: 'cancel'): () => void
}>();
const greeting = ref(props.greeting);

function onDelete() {
  console.log(greeting);
  emits('cancel');
}
</script>

<template>
  <div>
    <v-list-item>
      <div class="d-flex justify-space-between">
        <h1>Edit variable</h1>
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
            v-model="greeting.userName" 
            label="Name" 
            :rules="[
              v => !!v || 'Field is required'
            ]"
          />

          <v-text-field 
            v-model="greeting.text" 
            label="Message for sending in chat" 
          />

          <v-checkbox 
            v-model="greeting.isReply" 
            label="Use twitch reply feature" 
            class="mb-2" 
            density="compact"
          />
        </div>
      </v-form>
    </v-list-item>
  </div>
</template>