<script lang="ts" setup>
import { mdiPencil, mdiTrashCan } from '@mdi/js';
import { ChannelTimer } from '@tsuwari/typeorm/entities/ChannelTimer';
import { computed, ref } from 'vue';
import { useDisplay } from 'vuetify/lib/framework.mjs';

import confirmDeletion from '@/components/confirmDeletion.vue';
import TimersDrawer from '@/components/drawers/timerEdit.vue';
import { timers } from '@/data/timers';

const display = useDisplay();
const isMobile = computed(() => {
  return display.xs.value || display.sm.value;
});

const isTimerEdit = ref(false);
const editableTimer = ref<ChannelTimer | undefined>();

function setEditTimer(c: ChannelTimer) {
  isTimerEdit.value = true;
  editableTimer.value = c;
}

function cancelEdit() {
  isTimerEdit.value = false;
  editableTimer.value = undefined;
}

function onDelete(t: ChannelTimer) {
  console.log(t);
}
</script>

<template>
  <div class="d-flex flex-column">
    <v-table class="mt-2">
      <thead>
        <tr>
          <th class="text-left">
            Name
          </th>
          <th class="text-left">
            Responses
          </th>
          <th class="text-left">
            Status
          </th>
          <th class="text-left">
            Actions
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="timer in timers"
          :key="timer.name"
        >
          <td>
            {{ timer.name }}
          </td>
          <td v-if="!isMobile">
            <div v-if="timer.responses?.length > 1">
              <p 
                v-for="(r, i) of timer.responses"
               
                :key="i"
              >
                {{ r.text }}
                <br />
              </p>
            </div>
            
            <p v-else>
              {{ timer.responses[0]!.text ?? '' }}
            </p>
          </td>
          <td>
            <v-switch v-model="timer.enabled" class="d-flex justify-center" color="green"></v-switch>
          </td>
          <td>
            <div class="d-flex flex-row">
              <v-btn :icon="mdiPencil" size="x-small" color="purple" @click="setEditTimer(timer)" />
              <confirmDeletion :cb="() => onDelete(timer)">
                <v-btn :icon="mdiTrashCan" size="x-small" color="red" class="ml-2" @click="$emit('click')" />
              </confirmDeletion>
            </div>
          </td>
        </tr>
      </tbody>
    </v-table>


    <v-navigation-drawer
      v-if="isTimerEdit"
      v-model="isTimerEdit"
      temporary
      location="right"
      :class="[isMobile ? 'w-100' : 'w-50']"
    >
      <TimersDrawer :timer="editableTimer!" @cancel="cancelEdit" />
    </v-navigation-drawer>
  </div>
</template>