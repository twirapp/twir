<script lang="ts" setup>
import { mdiPencil, mdiTrashCan } from '@mdi/js';
import { ChannelKeyword } from '@tsuwari/typeorm/entities/ChannelKeyword';
import { computed, ref } from 'vue';
import { useDisplay } from 'vuetify';

import confirmDeletion from '@/components/confirmDeletion.vue';
import KeywordDrawer from '@/components/drawers/keywordEdit.vue';
import { keywords } from '@/stores/keywords';

const display = useDisplay();
const isMobile = computed(() => {
  return display.xs.value || display.sm.value;
});

const isKeywordEdit = ref(false);
const editableKeyword = ref<ChannelKeyword | undefined>();

function setEditKeyword(c: ChannelKeyword) {
  isKeywordEdit.value = true;
  editableKeyword.value = c;
}

function cancelEdit() {
  isKeywordEdit.value = false;
  editableKeyword.value = undefined;
}

function onDelete(k: ChannelKeyword) {
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
            Responses
          </th>
          <th class="text-left">
            Usages
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
          v-for="keyword in keywords"
          :key="keyword.id"
        >
          <td>
            {{ keyword.text }}
          </td>
          <td>
            <p v-if="keyword.response">
              {{ keyword.response }}
            </p>
            <v-chip v-else>
              No response
            </v-chip>
          </td>
          <td>
            {{ keyword.usages }}
          </td>
          <td>
            <v-switch v-model="keyword.enabled" class="d-flex justify-center" color="green"></v-switch>
          </td>
          <td>
            <div class="d-flex flex-row">
              <v-btn :icon="mdiPencil" size="x-small" color="purple" @click="setEditKeyword(keyword)" />
              <confirmDeletion :cb="() => onDelete(keyword)">
                <v-btn :icon="mdiTrashCan" size="x-small" color="red" class="ml-2" />
              </confirmDeletion>
            </div>
          </td>
        </tr>
      </tbody>
    </v-table>


    <v-navigation-drawer
      v-if="isKeywordEdit"
      v-model="isKeywordEdit"
      temporary
      location="right"
      :class="[isMobile ? 'w-100' : 'w-50']"
    >
      <KeywordDrawer :keyword="editableKeyword!" @cancel="cancelEdit" />
    </v-navigation-drawer>
  </div>
</template>