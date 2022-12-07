<script lang="ts" setup>
import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import { ref } from 'vue';

import confirmDeletion from '../confirmDeletion.vue';

const props = defineProps<{
  keyword: ChannelKeyword
}>();

const emits = defineEmits<{
  (event: 'cancel'): () => void
}>();

const keyword = ref(props.keyword);

function onDelete() {
  console.log(keyword);
  emits('cancel');
}
</script>

<template>
  <div>
    <v-list-item>
      <div class="d-flex justify-space-between">
        <h1>Edit keyword</h1>
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
            v-model="keyword.text" 
            label="Text" 
            :rules="[
              v => !!v || 'Field is required'
            ]"
          />

          <div>
            <h4>Cooldown (seconds)</h4>
            <div class="d-flex justify-space-between mt-2">
              <v-slider
                v-model="keyword.cooldown"
                :min="1"
                :max="120"
                :step="1"
              >
                <template #append>
                  <v-text-field
                    v-model="keyword.cooldown"
                    hide-details
                    single-line
                    variant="outlined"
                    density="compact"
                    type="number"
                    style="width: 100px"
                  ></v-text-field>
                </template>
              </v-slider>
            </div>
          </div>

          <v-textarea
            v-model="keyword.response"
            auto-grow
            rows="1"
            row-height="5"
            class="mt-2"
            label="Response"
          >
          </v-textarea>

          <v-text-field
            v-model="keyword.usages"
            label="Usages"
            type="number"
          ></v-text-field>
        </div>
      </v-form>
    </v-list-item>
  </div>
</template>