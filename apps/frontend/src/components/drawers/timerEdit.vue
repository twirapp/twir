<script lang="ts" setup>
import { mdiPlus, mdiClose } from '@mdi/js';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import { ref } from 'vue';

import confirmDeletion from '@/components/confirmDeletion.vue';

const props = defineProps<{
  timer: ChannelTimer
}>();

const emits = defineEmits<{
  (event: 'cancel'): () => void
}>();

const timer = ref(props.timer);

function onDelete() {
  console.log(timer);
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
            v-model="timer.name" 
            label="Name" 
            :rules="[
              v => !!v || 'Field is required'
            ]"
          />

          <div>
            <h4>Interval (minutes)</h4>
            <div class="d-flex justify-space-between mt-2">
              <v-slider
                v-model="timer.timeInterval"
                :min="1"
                :max="120"
                :step="1"
              >
                <template #append>
                  <v-text-field
                    v-model="timer.timeInterval"
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

          <div>
            <h4>Interval (messages)</h4>
            <div class="d-flex justify-space-between mt-2">
              <v-slider
                v-model="timer.messageInterval"
                :min="1"
                :max="120"
                :step="1"
              >
                <template #append>
                  <v-text-field
                    v-model="timer.messageInterval"
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

          <v-divider></v-divider>
          <div class="d-flex justify-space-between mt-2">
            <h4>Responses</h4>
            <div>
              <v-btn 
                variant="outlined" 
                size="x-small" 
                @click="timer!.responses!.push({
                  order: timer!.responses?.length ? timer!.responses?.length - 1 : 0,
                  text: '',
                } as any)"
              >
                <v-icon>{{ mdiPlus }}</v-icon>
              </v-btn>
            </div>
          </div>

          <v-sheet
            v-if="!timer.responses?.length"
            rounded
            class="mt-2 pa-4"
            color="#484749"
          >
            No responses added
          </v-sheet>

          <v-textarea
            v-for="(response, index) of timer!.responses" 
            :key="index"
            v-model="timer!.responses![index].text"
            auto-grow
            :append-icon="mdiClose"
            rows="1"
            row-height="5"
            class="mt-2"
            :rules="[
              v => !!v || 'Field is required'
            ]"
            @click:append="timer.responses!.splice(index, 1)"
          >
            <template #prepend>
              <v-checkbox v-model="timer.responses[index].isAnnounce" label="Use announce"></v-checkbox>
            </template>
          </v-textarea>
        </div>
      </v-form>
    </v-list-item>
  </div>
</template>