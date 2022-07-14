<script lang="ts" setup>
export type VariablesList = Array<{ name: string, example?: string, description?: string }>

import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { UpdateOrCreateCommandDto } from '@tsuwari/api/src/v1/commands/dto/create';
import { computed, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import Command from '../components/Command.vue';
import { VariableType } from './Variables.vue';

import Add from '@/assets/buttons/add.svg';
import Dota2Icon from '@/assets/icons/dota2.svg?component';
import Integrations from '@/assets/sidebar/integrations.svg?component';
import MyBtn from '@/components/elements/MyBtn.vue';
import { useUpdatingData } from '@/functions/useUpdatingData';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

type CommandType = UpdateOrCreateCommandDto & { new?: boolean, default?: boolean, defaultName?: string }

const { data: axiosData } = useUpdatingData(`/v1/channels/{dashboardId}/commands`);

const commands = ref<CommandType[]>([]);
const variablesList = ref<VariablesList>([]);
const currentEditableCommand = ref<CommandType | null>(null);
const searchFilter = ref<string>('');
const filteredCommands = computed(() => {
  return commands.value.filter(c => c.name).filter(c => searchFilter.value ? [c.name, ...c.aliases as string[]].some(s => s.includes(searchFilter.value)) : true).sort((a, b) => a.name.localeCompare(b.name));
});
const { t } = useI18n({
  useScope: 'global',
});

const DotaGroup = new Set(['np', 'dota addacc', 'dota delacc', 'wl', 'dota listacc']);

watch(axiosData, (v: CommandType[]) => {
  commands.value = v;
  currentEditableCommand.value = filteredCommands.value[0];
});

selectedDashboardStore.subscribe(async (v) => {
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
  <div class="block flex justify-between md:hidden mx-2 my-2 space-x-2">
    <Popover
      class="relative"
    >
      <PopoverButton class="bg-green-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-green-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition">
        <svg
          fill="none"
          width="16"
          height="16"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M4 6h16M4 12h16m-7 6h7"
        /></svg>
      </PopoverButton>

      <PopoverPanel
        v-slot="{ close }"
        :focus="true"
        class="absolute bg-[#121010] z-10"
      >
        <ul class="max-h-[55vh] overflow-y-auto scrollbar scrollbar-thin scrollbar-thumb-gray-900 scrollbar-track-gray-600">
          <li
            v-for="command, index of filteredCommands"
            :key="index"
            class="px-0.5"
            :class="{ 'border-l-2': filteredCommands.indexOf(currentEditableCommand!) === index }"
            @click="() => {
              if (!currentEditableCommand!.id) commands.splice(commands.indexOf(currentEditableCommand!), 1)
              currentEditableCommand = command
              close()
            }"
          >
            <button
              aria-current="page"
              href="/dashboard/commands"
              class="border-slate-300 duration-300 ease-in-out flex h-8 hover:bg-[#202122] items-center justify-between mt-0 overflow-hidden px-2 ripple-surface-primary text-ellipsis text-sm text-white transition w-full whitespace-nowrap"
              :class="{
                'bg-neutral-700': filteredCommands.indexOf(currentEditableCommand!) === index
              }"
            > 
              <span>{{ command.name }}</span>
              <Dota2Icon
                v-if="command.defaultName && DotaGroup.has(command.defaultName)"
                class="h-[17px] ml-3"
              />
            <!-- <Integrations /> -->
            </button>
          </li>
        </ul>
      </PopoverPanel>
    </Popover>

    <div>
      <MyBtn
        color="green"
        size="default"
        @click="insertCommand"
      >
        <Add />
      </MyBtn>
    </div>
  </div>

  <div class="flex h-full">
    <div class="border-gray-700 border-r hidden md:block rounded w-40">
      <button
        class="bg-green-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium grid hover:bg-green-700 inline-block leading-tight m-auto place-items-center px-6 py-2.5 shadow text-white text-xs transition uppercase w-full"
        @click="insertCommand"
      >
        <Add />
      </button>
      <div class="form-floating">
        <input
          id="searchCommand"
          v-model="searchFilter"
          type="text"
          class="bg-clip-padding
                    bg-white
                    border
                    border-gray-300
                    border-solid
                    ease-in-out
                    focus:outline-none
                    font-normal
                    form-control
                    text-base
                    text-gray-700
                    transition
                    w-full"
          placeholder="command"
        >
        <label
          for="searchCommand"
          class="text-gray-700"
        >{{ t('pages.commands.searchCommand') }}</label>
      </div>


      <ul class="max-h-[75vh] menu overflow-auto scrollbar scrollbar-thin scrollbar-thumb-gray-900 scrollbar-track-gray-600">
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
            class="border-slate-300 duration-300 ease-in-out flex h-8 hover:bg-[#202122] items-center justify-between mt-0 overflow-hidden px-2 ripple-surface-primary text-ellipsis text-sm text-white transition w-full whitespace-nowrap"
            :class="{
              'bg-neutral-700': filteredCommands.indexOf(currentEditableCommand!) === index
            }"
          > 
            <span>{{ command.name }}</span>
            <Dota2Icon
              v-if="command.defaultName && DotaGroup.has(command.defaultName)"
              class="h-[17px]"
            />
          </button>
        </li>
      </ul>
    </div>

    <div
      v-if="currentEditableCommand"
      class="card h-fit m-1.5 max-w-2xl md:m-3 p-1 rounded shadow sm:block text-white w-full"
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
