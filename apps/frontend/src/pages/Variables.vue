<script lang="ts" setup>
import { mdiPencil, mdiTrashCan } from '@mdi/js';
import { type ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { ref } from 'vue';
import { useDisplay } from 'vuetify';

import confirmDeletion from '@/components/confirmDeletion.vue';
import KeywordDrawer from '@/components/drawers/keywordEdit.vue';
import { variables } from '@/data/variables';

const { mobile } = useDisplay();

const isVariableEdit = ref(false);
const editableVariable = ref<ChannelCustomvar | undefined>();

function setEditKeyword(c: ChannelCustomvar) {
  isVariableEdit.value = true;
  editableVariable.value = JSON.parse(JSON.stringify(c));
}

function cancelEdit() {
  isVariableEdit.value = false;
  editableVariable.value = undefined;
}

function deleteVariable(v: ChannelCustomvar) {

}

function onDelete(v: ChannelCustomvar) {
  console.log(v);
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
            Response
          </th>
          <th class="text-left">
            Type
          </th>
          <th class="text-left">
            Response
          </th>
          <th class="text-left">
            Actions
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="variable in variables"
          :key="variable.id"
        >
          <td>{{ variable.name }}</td>
          <td>{{ variable.response }}</td>
          <td>
            <v-chip size="small">
              {{ variable.type }}
            </v-chip>
          </td>
          <td>
            <p v-if="variable.response">
              {{ variable.response }}
            </p>
            <p v-else>
              <v-chip size="small" color="secondary">
                Cannot display script
              </v-chip>
            </p>
          </td>
          <td>
            <div class="d-flex flex-row">
              <v-btn :icon="mdiPencil" size="x-small" color="purple" @click="setEditKeyword(variable)" />
              <confirmDeletion :cb="() => deleteVariable(variable)">
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
      <KeywordDrawer :keyword="editableVariable!" @cancel="cancelEdit" />
    </v-navigation-drawer>
  </div>
</template>