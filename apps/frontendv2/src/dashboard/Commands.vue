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
        <div class="p-4">
          <!-- <div class="flex">
            <p>Command status</p>
            <input
              v-model="command.enabled"
              type="checkbox"
              class="toggle"
              checked
            >
          </div> -->

          <div
            class="grid grid-cols-2 gap-1"
          >
            <div>
              <div class="label">
                <span class="label-text">Name</span>
              </div>
              <input
                v-model="command.name"
                type="text"
                placeholder="uptime"
                :disabled="!command.edit"
                class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
              >
            </div>


            <div>
              <span class="label text-center">Cooldown</span>
            
              <div class="grid grid-cols-2">
                <input
                  v-model="command.cooldown"
                  type="number"
                  placeholder="0"
                  :disabled="!command.edit"
                  class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-4/5 input-sm"
                >
                
                <select
                  v-model="command.cooldownType"
                  :disabled="!command.edit"
                  class="form-control px-3 py-1.5 text-gray-700 rounded select select-sm w-full"
                >
                  <option
                    v-for="type of Object.entries(cooldownType)"
                    :key="type[0]"
                    :value="type[1]"
                  >
                    {{ type[0] }}
                  </option>
                </select>
              </div>
            </div>

            <div>
              <span class="label">Permission</span>
              <select
                v-model="command.permission"
                :disabled="!command.edit"
                class="form-control px-3 py-1.5 text-gray-700 rounded select select-sm w-full mb-2"
              >
                <option
                  disabled
                  selected
                >
                  Choose permission
                </option>
                <option
                  v-for="permission of Object.entries(perms)"
                  :key="permission[0]"
                  :value="permission[1]"
                >
                  {{ permission[0] }}
                </option>
              </select>
            </div>

            <div>
              <span class="label">Description (optional)</span>
              <input
                v-model="command.description"
                type="text"
                :disabled="!command.edit"
                placeholder="great command ;)"
                class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
              >
            </div>

            <div class="col-span-2">
              <span class="label">
                <span>Responses
                  <button
                    v-if="command.edit"
                    class="px-2 py-0.5 inline-block bg-green-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:shadow-lg focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out"
                    @click="command.responses.push({ text: '' })"
                  >
                    +
                  </button>
                </span>
              </span>

              <div class="input-group grid grid-cols-1 pt-1 gap-2 max-h-40 scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600">
                <div
                  v-for="_response, responseIndex in command.responses"
                  :key="responseIndex"
                  class="flex flex-wrap items-stretch mb-4 relative"
                  style="width: 99%;"
                >
                  <input
                    v-model.lazy="command.responses[responseIndex].text"
                    type="text"
                    :disabled="!command.edit"
                    class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1 border border-grey-light text-gray-700 rounded rounded-r-none px-3 py-1.5 relative"
                    placeholder="command response"
                  >
                  <div
                    v-if="command.edit"
                    class="flex -mr-px cursor-pointer"
                    @click="command.responses?.splice(responseIndex, 1)"
                  >
                    <span class="flex items-center leading-normal bg-red-500 rounded rounded-l-none border-0 border-l-0 border-grey-light px-5 py-1.5 whitespace-no-wrap text-grey-dark text-sm">X</span>
                  </div>
                </div>
              </div>
            </div>

            <div class="col-span-2">
              <span class="label">  
                <span>Aliases
                  <button
                    v-if="command.edit"
                    class="px-2 py-0.5 inline-block bg-green-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md  hover:shadow-lg focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out"
                    @click="command.aliases?.push('')"
                  >
                    +
                  </button>
                </span>
              </span>

              <div class="input-group pt-1 pr-2 grid lg:grid-cols-2 md:grid-cols-2 sm:grid-cols-2 grid-cols-2 xl:grid-cols-3 gap-1 max-h-40 scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600">
                <div
                  v-for="aliase, aliaseIndex in command.aliases"
                  :key="aliase"
                  class="flex flex-wrap items-stretch mb-4 relative"
                >
                  <input
                    v-model.lazy="command.aliases![aliaseIndex]"
                    type="text"
                    class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1 border border-grey-light text-gray-700 rounded rounded-r-none px-3 py-1.5 relative"
                  >
                  <div
                    class="flex -mr-px cursor-pointer"
                    @click="command.aliases?.splice(aliaseIndex, 1)"
                  >
                    <span class="flex items-center leading-normal bg-red-500 rounded rounded-l-none border-0 border-l-0 border-grey-light px-5 py-1.5 whitespace-no-wrap text-grey-dark text-sm">X</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="card-actions flex justify-between mt-5">
            <div>
              <button
                v-if="!command.edit"
                type="button"
                class="inline-block px-6 py-2.5 bg-gray-200 text-gray-700 font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-gray-300 hover:shadow-lg focus:bg-gray-300 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-gray-400 active:shadow-lg transition duration-150 ease-in-out"
                @click="() => {
                  command.edit = true;
                  if (command.id) commandsBeforeEdit?.push(JSON.parse(JSON.stringify(command)))
                }"
              >
                Edit
              </button>
              <button
                v-else
                class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-purple-700 hover:shadow-lg focus:bg-purple-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out"
                @click="cancelEdit(command, commandIndex)"
              >
                Cancel
              </button>
            </div>
            <div v-if="command.edit">
              <button
                type="button"
                class="inline-block px-6 py-2.5 bg-red-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-red-700 hover:shadow-lg focus:bg-red-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-red-800 active:shadow-lg transition duration-150 ease-in-out"
                @click="deleteCommand(commandIndex, command.id)"
              >
                Delete
              </button>
              <button
                type="button"
                class="inline-block ml-2 px-6 py-2.5 bg-green-500 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-green-600 hover:shadow-lg focus:bg-green-600 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-green-700 active:shadow-lg transition duration-150 ease-in-out"
                @click="saveCommand(command, commandIndex)"
              >
                Save
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: linear-gradient(0deg, rgba(255, 255, 255, 0.05), rgba(255, 255, 255, 0.05)), #121212;
}
</style>