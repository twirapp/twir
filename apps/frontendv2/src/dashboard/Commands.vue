<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { UpdateOrCreateCommandDto } from '@tsuwari/api/src/v1/commands/dto/create';
import { CommandPermission, CooldownType } from '@tsuwari/prisma';
import { useTitle } from '@vueuse/core';
import { useAxios } from '@vueuse/integrations/useAxios';
import { ref, watch } from 'vue';

import Command from '../components/Command.vue';

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

type CommandType = UpdateOrCreateCommandDto & { 
  edit?: boolean
}

const selectedDashboard = useStore(selectedDashboardStore);

const { execute, data: axiosData, isLoading } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/commands`, api, { immediate: false });

const commands = ref<CommandType[]>([]);
const commandsBeforeEdit = ref<CommandType[]>([]);

watch(axiosData, (v: CommandType[]) => {
  commands.value = v;
  commandsBeforeEdit.value = [];
});

selectedDashboardStore.subscribe((v) => {
  execute(`/v1/channels/${v.channelId}/commands`);
});

function insertCommand() {
  if (commands.value) {
    const command: CommandType = {
      name: '',
      aliases: [],
      cooldown: 0,
      permission: 'VIEWER',
      description: null,
      visible: true,
      enabled: true,
      responses: [
        { text: '' },
      ],
      edit: true,
      cooldownType: 'GLOBAL',
    };

    commands.value.unshift(command);
  }
}
</script>

<template>
  <div class="p-1">
    <div class="flow-root">
      <div class="float-left rounded btn btn-primary btn-sm w-full mb-1 md:w-auto">
        <button
          class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-purple-700 hover:shadow-lg focus:bg-purple-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out"
          @click="insertCommand"
        >
          Add new command
        </button>
      </div>

      <!-- <input
        type="text"
        placeholder="Search by keyword..."
        class="float-right rounded input input-sm input-bordered w-full md:w-60"
      > -->
    </div>
  </div>

  <div class="w-full">
    <div v-if="isLoading">
      <div class="flex items-center justify-center ">
        <svg
          class="animate-spin -ml-1 mr-3 h-24 w-24 text-white"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
        >
          <circle
            class="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            stroke-width="4"
          />
          <path
            class="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          />
        </svg>
      </div>
    </div>
    <div 
      v-else
      class="grid lg:grid-cols-3 md:grid-cols-2 grid-cols-1 gap-2"
    >
      <div
        v-for="command, commandIndex in commands"
        :key="commandIndex"
        class="block rounded-lg card text-white shadow-lg"
      >
        <Command
          :command="command"
          :commands="commands"
          :commands-before-edit="commandsBeforeEdit"
        />
      </div>
    </div>
  </div>
</template>

<style>
.card {
  background: linear-gradient(0deg, rgba(255, 255, 255, 0.05), rgba(255, 255, 255, 0.05)), #121212;
}
</style>