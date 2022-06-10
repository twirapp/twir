<script lang="ts" setup>
export type VariablesList = Array<{ name: string, example?: string, description?: string }>

import { useStore } from '@nanostores/vue';
import { UpdateOrCreateCommandDto } from '@tsuwari/api/src/v1/commands/dto/create';
import { useTitle } from '@vueuse/core';
import { useAxios } from '@vueuse/integrations/useAxios';
import { computed, ref, watch } from 'vue';

import Command from '../components/Command.vue';
import { VariableType } from './Variables.vue';

import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const title = useTitle();
title.value = 'Tsuwari - Commands';

type CommandType = UpdateOrCreateCommandDto & { new?: boolean, default?: boolean | null }

const selectedDashboard = useStore(selectedDashboardStore);

const { execute, data: axiosData } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/commands`, api, { immediate: false });

const commands = ref<CommandType[]>([]);
const variablesList = ref<VariablesList>([]);
const currentEditableCommand = ref<CommandType | null>(null);
const searchFilter = ref<string>('');
const filteredCommands = computed(() => {
  return commands.value.filter(c => c.name).filter(c => searchFilter.value ? [c.name, ...c.aliases as string[]].some(s => s.includes(searchFilter.value)) : true).sort((a, b) => a.name.localeCompare(b.name));
});

watch(axiosData, (v: CommandType[]) => {
  commands.value = v;
  currentEditableCommand.value = filteredCommands.value[0];
});

selectedDashboardStore.subscribe(async (v) => {
  execute(`/v1/channels/${v.channelId}/commands`);
  const [custom, builtIn] = await Promise.all([
    api(`v1/channels/${v.channelId}/variables`),
    api(`v1/channels/${v.channelId}/variables/builtin`),
  ]);

  variablesList.value = [
    ...custom.data.map((c: VariableType) => ({ name: c.name, example: `customvar|${c.name}`, description: `Created custom variable ${c.name.toUpperCase()}` })),
    ...builtIn.data,
  ];
});

function insertCommand() {
  if (commands.value && !currentEditableCommand.value?.new) {
    const command: CommandType = {
      name: '',
      aliases: [],
      cooldown: 5,
      permission: 'VIEWER',
      description: null,
      visible: true,
      enabled: true,
      responses: [],
      cooldownType: 'GLOBAL',
      new: true,
    };

    commands.value.unshift(command);
    currentEditableCommand.value = command;
  }
}


function deleteCommand(index: number) {
  commands.value = commands.value.filter((_, i) => i !== index);
  currentEditableCommand.value = commands.value[0];
}

function onSave(index: number) {
  currentEditableCommand.value = commands.value[index];
}
</script>

<template>
  <div class="flex">
    <div>
      <div class="w-40 h-[90%] rounded border-r border-b border-gray-700">
        <button
          class="px-6 py-2.5 w-full inline-block bg-green-500 text-white font-medium text-xs leading-tight uppercase shadow-md hover:bg-green-500 hover:shadow-lg focus:bg-green-600 focus:shadow-lg focus:outline-none focus:ring-0 active:shadow-lg transition duration-150 ease-in-out"
          @click="insertCommand"
        >
          +
        </button>
        <div class="form-floating">
          <input
            id="searchCommand"
            v-model="searchFilter"
            type="text"
            class="form-control
                    w-full
                    text-base
                    font-normal
                    text-gray-700
                    bg-white bg-clip-padding
                    border border-solid border-gray-300
                    transition
                    ease-in-out
                    m-0
                    focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none"
            placeholder="command"
          >
          <label
            for="searchCommand"
            class="text-gray-700"
          >Search command</label>
        </div>
        <ul class="menu max-h-screen min-h-screen scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600">
          <li
            v-for="command, index of filteredCommands
            "
            :key="index"
            :class="{ 'border-l-2': filteredCommands.indexOf(currentEditableCommand!) === index }"
            @click="() => {
              if (!currentEditableCommand!.id) commands.splice(commands.indexOf(currentEditableCommand!), 1)
              currentEditableCommand = command  
            }"
          >
            <button
              aria-current="page"
              href="/dashboard/commands"
              class="flex items-center mt-0 text-sm px-2 h-8 w-full overflow-hidden text-white text-ellipsis whitespace-nowrap hover:bg-[#202122] border-slate-300 transition duration-300 ease-in-out ripple-surface-primary"
              :class="{
                'bg-neutral-700': filteredCommands.indexOf(currentEditableCommand!) === index
              }"
            >
              <span class="w-3 h-3" /><span>{{ command.name }}</span>
            </button>
          </li>
        </ul>
      </div>
    </div>

    <div
      v-if="currentEditableCommand"
      class="w-full p-1 hidden sm:block h-fit m-4 block max-w-2xl rounded-lg card text-white shadow-lg"
    >
      <Command 
        :command="currentEditableCommand" 
        :commands="commands" 
        :variables-list="variablesList"
        @delete="deleteCommand"
        @save="onSave"
      />
    </div>
  </div>
</template>
