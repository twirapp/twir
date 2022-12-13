<script lang="ts" setup>
import { mdiPencil, mdiTrashCan } from '@mdi/js';
import { type ChannelCustomvar } from '@tsuwari/typeorm/entities/ChannelCustomvar';
import { computed, ref } from 'vue';
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

function onDelete(k: ChannelCustomvar) {
  console.log(k);
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
            Type
          </th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="variable in variables"
          :key="variable.id"
        >
          <td>{{ variable.response }}</td>
          <td>
            <v-chip size="small">
              {{ variable.type }}
            </v-chip>
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