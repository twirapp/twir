<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { UpdateOrCreateCommandDto } from '@tsuwari/api/src/v1/commands/dto/create';
import { CommandPermission, CooldownType } from '@tsuwari/prisma';
import { useTitle } from '@vueuse/core';
import { useAxios } from '@vueuse/integrations/useAxios';
import { ref, watch } from 'vue';

import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';



const title = useTitle();
title.value = 'Tsuwari - Commands';

const perms = {
  'Broadcaster': 'BROADCASTER',
  'Moderator\'s': 'MODERATOR',
  'Vip\'s': 'VIP',
  'Subscriber\'s': 'SUBSCRIBER',
  'Viewers': 'VIEWER',
} as { [x: string]: CommandPermission };

const cooldownType = {
  'Global': 'GLOBAL',
  'Per user': 'PER_USER',
} as { [x: string]: CooldownType };

type CommandType = UpdateOrCreateCommandDto & { edit?: boolean }

const selectedDashboard = useStore(selectedDashboardStore);

const { execute, data: axiosData, isLoading } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/commands`, api, { immediate: false });

const commands = ref<CommandType[]>();
const commandsBeforeEdit = ref<CommandType[]>();

watch(axiosData, (v) => {
  commands.value = v;
  commandsBeforeEdit.value = [];
});

selectedDashboardStore.subscribe((v) => {
  execute(`/v1/channels/${v.channelId}/commands`);
  // setCommands();
});

/* async function setCommands() {
  const { data } = await api(`/v1/channels/${selectedDashboard.value.channelId}/commands`);
  commands.value = data;
} */

async function deleteCommand(index: number, id?: string) {
  if (id) {
    await api.delete(`/v1/channels/${selectedDashboard.value.channelId}/commands/${id}`);
  }

  if (commands.value) {
    commands.value = commands.value.filter((_, i) => i !== index);
  }
}

function insertCommand() {
  if (commands.value) {
    commands.value.unshift({
      name: '',
      aliases: [],
      cooldown: 0,
      permission: 'VIEWER',
      description: null,
      visible: true,
      enabled: true,
      responses: [
        { text: null },
      ],
      edit: true,
      cooldownType: 'GLOBAL',
    });
  }
}

async function saveCommand(command: CommandType, index: number) {
  let data;
  if (command.id) {
    const request = await api.put(`/v1/channels/${selectedDashboard.value.channelId}/commands/${command.id}`, command);  
    data = request.data;
  } else {
    const request = await api.post(`/v1/channels/${selectedDashboard.value.channelId}/commands`, command);
    data = request.data;
  }

  if (commands.value && commands.value[index]) {
    commands.value[index] = data;
  }
}

function cancelEdit(command: CommandType, index: number) {
  if (command.id && commands.value) {
    const editableCommand = commandsBeforeEdit.value?.find(c => c.id === command.id);
    if (editableCommand) {
      commands.value[index] = {
        ...editableCommand,
        edit: false,
      };
      commandsBeforeEdit.value?.splice(commandsBeforeEdit.value.indexOf(editableCommand), 1);
    }
  } else {
    commands.value?.splice(index, 1);
  }
}

</script>

<template>
  <div class="q-pa-md">
    <div class="row">
      <q-card
        flat
        bordered
        class="col-4"
      >
        <q-card-section>
          <div class="text-h6">
            Status
          </div>
        </q-card-section>
      </q-card>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: linear-gradient(0deg, rgba(255, 255, 255, 0.05), rgba(255, 255, 255, 0.05)), #121212;
}
</style>