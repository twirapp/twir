<script lang="ts" setup>
import { mdiPencil, mdiTrashCan } from '@mdi/js';
import { ref } from 'vue';
import { useDisplay } from 'vuetify';

import confirmDeletion from '@/components/confirmDeletion.vue';
import VariableDrawer from '@/components/drawers/variableEdit.vue';
import { greetings, type Greeting } from '@/stores/greetings';

const { mobile } = useDisplay();
const isGreetingEdit = ref(false);
const editableGreeting = ref<Greeting | undefined>();

function setEditKeyword(c: Greeting) {
  isGreetingEdit.value = true;
  editableGreeting.value = JSON.parse(JSON.stringify(c));
}

function cancelEdit() {
  isGreetingEdit.value = false;
  editableGreeting.value = undefined;
}

function deleteVariable(g: Greeting) {}

function onDelete(g: Greeting) {
  console.log(g);
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
            Text
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
          v-for="greeting in greetings"
          :key="greeting.id"
        >
          <td>{{ greeting.userName }}</td>
          <td>{{ greeting.text }}</td>
          <td><v-switch v-model="greeting.enabled" class="d-flex justify-center" color="green"></v-switch></td>
          <td>
            <div class="d-flex flex-row">
              <v-btn :icon="mdiPencil" size="x-small" color="purple" @click="setEditKeyword(greeting)" />
              <confirmDeletion :cb="() => deleteVariable(greeting)">
                <v-btn :icon="mdiTrashCan" size="x-small" color="red" class="ml-2" @click="$emit('click')" />
              </confirmDeletion>
            </div>
          </td>
        </tr>
      </tbody>
    </v-table>


    <v-navigation-drawer
      v-if="isVariableEdit"
      v-model="isVariableEdit"
      temporary
      location="right"
      :class="[mobile ? 'w-100' : 'w-50']"
    >
      <VariableDrawer :variable="editableVariable!" @cancel="cancelEdit" />
    </v-navigation-drawer>
  </div>
</template>