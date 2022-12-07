<script lang="ts" setup>
import { mdiPencil, mdiTrashCan } from '@mdi/js';
import { ChannelCommand, type CommandModule } from '@tsuwari/typeorm/entities/ChannelCommand';
import { computed, ref } from 'vue';
import { useDisplay } from 'vuetify/lib/framework.mjs';

import confirmDeletion from '@/components/confirmDeletion.vue';
import CommandDrawer from '@/components/drawers/commandEdit.vue';
import { commands } from '@/data/commands';
import { tabs } from '@/data/commandsTabs';

const selectedTab = ref<keyof typeof CommandModule>('CUSTOM');

const { smAndDown } = useDisplay();

const commandsList = computed(() => {
  return commands.value.filter(c => c.module === selectedTab.value);
});

const isCommandEdit = ref(false);
const editableCommand = ref<ChannelCommand | undefined>();

function setEditCommand(c: ChannelCommand) {
  isCommandEdit.value = true;
  editableCommand.value = c;
}

const symbolsRegexp = /^\W$|_|/;
function preventSymbolsInCommandName(e: KeyboardEvent) {
  // if (symbolsRegexp.test(e.key)) {
  //   e.preventDefault();
  // }
}

const c = console;

function cancelEdit() {
  isCommandEdit.value = false;
  editableCommand.value = undefined;
}

function deleteCommand(id: string) {
  return null;
}
</script>

<template>
  <div class="d-flex flex-column">
    <div>
      <v-tabs
        v-model="selectedTab"
        centered
        stacked
        fixed-tabs
        class="d-flex justify-center"
      >
        <v-tab v-for="tab of tabs" :key="tab.value" :value="tab.value">
          <v-icon v-if="tab.icon">
            {{ tab.icon }}
          </v-icon>
          {{ tab.name }}
        </v-tab>
      </v-tabs>
    </div>
    

    <v-table class="mt-2">
      <thead>
        <tr>
          <th class="text-left">
            Name
          </th>
          <th v-if="!smAndDown && selectedTab === 'CUSTOM'" class="text-left">
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
          v-for="command in commandsList"
          :key="command.name"
        >
          <td :width="(!smAndDown && selectedTab === 'CUSTOM') ? '10%' : '100%'">
            {{ command.name }}
          </td>
          <td v-if="!smAndDown && selectedTab === 'CUSTOM'">
            <p v-if="command.responses!.length > 1" v-html="command.responses?.map((t) => t.text).join('<br>')"></p>
            <p v-else>
              {{ command.responses![0]?.text ?? '' }}
            </p>
          </td>
          <td>
            <v-switch v-model="command.enabled" class="d-flex justify-center" color="green"></v-switch>
          </td>
          <td>
            <div class="d-flex flex-row">
              <v-btn :icon="mdiPencil" size="x-small" color="purple" @click="setEditCommand(command)" />
              <confirmDeletion @on-confirmed="() => deleteCommand(command.id)">
                <v-btn :icon="mdiTrashCan" size="x-small" color="red" class="ml-2" @click="$emit('click')" />
              </confirmDeletion>
            </div>
          </td>
        </tr>
      </tbody>
    </v-table>


    <v-navigation-drawer
      v-if="isCommandEdit"
      v-model="isCommandEdit"
      temporary
      location="right"
      :class="[smAndDown ? 'w-100' : 'w-50']"
    >
      <CommandDrawer :command="editableCommand!" @cancel="cancelEdit" />
    </v-navigation-drawer>
  </div>
</template>