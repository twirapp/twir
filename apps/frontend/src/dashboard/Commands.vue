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
  commandsBeforeEdit.value = v;
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

</script>

<template>
  <div class="p-1">
    <div class="flow-root">
      <div class="float-left rounded btn btn-primary btn-sm w-full mb-1 md:w-auto">
        <label>
          <svg
            width="20"
            height="20"
            viewBox="0 0 20 20"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M10 4.16663V15.8333"
              stroke="white"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            <path
              d="M4.16663 10H15.8333"
              stroke="white"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </label>

        <label
          class="ml-1 text-white"
          @click="insertCommand"
        >Add new command</label>
      </div>

      <input
        type="text"
        placeholder="Search by keyword..."
        class="float-right rounded input input-sm input-bordered w-full md:w-60"
      >
    </div>
  </div>

  <div class="w-full">
    <div v-if="isLoading">
      <div class="flex items-center justify-center ">
        <svg
          class="animate-spin -ml-1 h-24 w-24 text-white"
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
        class="card card-compact bg-base-200 drop-shadow-lg rounded"
      >
        <!-- <div v-if="command.edit" class="card-body grid grid-cols-2">
          <div>
            <div class="label">
              <span class="label-text">Name</span>
            </div>
            <input type="text" placeholder="uptime" v-model="command.name" class="rounded input input-bordered w-full input-sm" />
          </div>


          <div>
            <span class="label">Cooldown</span>
            
            <div class="grid grid-cols-2">
              <input type="number" placeholder="0" v-model="command.cooldown" class="rounded input input-bordered w-4/5 input-sm" />
               <select v-model="command.cooldownType" class="rounded select select-sm w-full mb-2">
                <option v-bind:key="type[0]" v-for="type of Object.entries(cooldownType)" :value="type[1]">
                  {{ type[0] }}
                </option>
              </select>
            </div>
          </div>

          <div>
            <span class="label">Permission</span>
            <select v-model="command.permission" class="rounded select select-sm w-full mb-2">
              <option disabled selected>Choose permission</option>
              <option v-bind:key="permission[0]" v-for="permission of Object.entries(perms)" :value="permission[1]">
                {{ permission[0] }}
              </option>
            </select>
          </div>

          <div>
            <span class="label">Description (optional)</span>
            <input type="text" placeholder="great command ;)" v-model="command.description" class="rounded input input-bordered w-full input-sm" />
          </div>

          <div class="col-span-2">
            <span class="label">
              <span>Responses <button @click="command.responses.push({ text: '' })" class="ml-3 btn btn-success btn-xs rounded">+</button></span>
            </span>
            <label class="input-group grid grid-cols-1 gap-2">
              <div v-bind:key="responseIndex" v-for="response, responseIndex in command.responses">
                <input type="text" placeholder="response of command" v-model.lazy="command.responses[responseIndex].text" class="input input-bordered input-sm w-5/6" />
                <button @click="command.responses?.splice(responseIndex, 1)" class="btn btn-error btn-sm">X</button>
              </div>
            </label>
          </div>

          <div class="col-span-2">
            <span class="label">
              <span>Aliases <button @click="command.aliases?.push('')" class="ml-3 rounded btn btn-success btn-xs">+</button></span>
            </span>
            <label class="input-group grid lg:grid-cols-2 md:grid-cols-2 sm:grid-cols-2 grid-cols-2 xl:grid-cols-3 gap-2">
              <div v-bind:key="aliase" v-for="aliase, aliaseIndex in command.aliases">
                <input type="text" v-model.lazy="command.aliases![aliaseIndex]" class="input input-bordered input-sm w-4/6" />
                <button @click="command.aliases?.splice(aliaseIndex, 1)" class="btn btn-error btn-sm">X</button>
              </div>
            </label>
          </div>

        </div>-->
        <!--<input v-model="command.enabled" type="checkbox" class="toggle" checked />-->
        <div
          v-if="command.edit"
          class="card-body"
        >
          <div class="flex">
            <p>Command status</p>
            <input
              v-model="command.enabled"
              type="checkbox"
              class="toggle"
              checked
            >
          </div>

          <div class="grid grid-cols-2 gap-1">
            <div>
              <div class="label">
                <span class="label-text">Name</span>
              </div>
              <input
                v-model="command.name"
                type="text"
                placeholder="uptime"
                class="rounded input input-bordered w-full input-sm"
              >
            </div>


            <div>
              <span class="label">Cooldown</span>
            
              <div class="grid grid-cols-2">
                <input
                  v-model="command.cooldown"
                  type="number"
                  placeholder="0"
                  class="rounded input input-bordered w-4/5 input-sm"
                >
                <select
                  v-model="command.cooldownType"
                  class="rounded select select-sm w-full mb-2"
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
                class="rounded select select-sm w-full mb-2"
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
                placeholder="great command ;)"
                class="rounded input input-bordered w-full input-sm"
              >
            </div>

            <div class="col-span-2">
              <span class="label">
                <span>Responses<button
                  class="ml-3 btn btn-success btn-xs rounded"
                  @click="command.responses.push({ text: '' })"
                >+</button></span>
              </span>

              <label class="input-group grid grid-cols-1 pt-1 gap-2 max-h-40 scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600">
                <div
                  v-for="response, responseIndex in command.responses"
                  :key="responseIndex"
                >
                  <label class="input-group pr-3">
                    <input
                      v-model.lazy="command.responses[responseIndex].text"
                      type="text"
                      placeholder="Response of command"
                      class="input input-bordered input-sm w-full"
                    >
                    <span
                      class="btn btn-error btn-sm no-animation"
                      @click="command.responses?.splice(responseIndex, 1)"
                    >-</span>
                  </label>
                </div>
              </label>
            </div>

            <div class="col-span-2">
              <span class="label">  
                <span>Aliases<button
                  class="ml-3 rounded btn btn-success btn-xs"
                  @click="command.aliases?.push('')"
                >+</button></span>
              </span>

              <label class="input-group pt-1 pr-2 grid lg:grid-cols-2 md:grid-cols-2 sm:grid-cols-2 grid-cols-2 xl:grid-cols-3 gap-1 max-h-40 scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600">
                <div
                  v-for="aliase, aliaseIndex in command.aliases"
                  :key="aliase"
                >
                  <input
                    v-model.lazy="command.aliases![aliaseIndex]"
                    type="text"
                    class="input input-bordered input-sm w-3/5"
                  >
                  <button
                    class="btn btn-error rounded btn-sm no-animation"
                    @click="command.aliases?.splice(aliaseIndex, 1)"
                  >X</button>
                </div>
              </label>
            </div>
          </div>
        </div>

        <div
          v-if="!command.edit"
          class="card-body p-4"
        >
          <div class="card-title">
            {{ command.name }}
          </div>
        </div>

        <div
          v-if="command.edit"
          class="card-actions flex justify-between m-3"
        >
          <div>
            <button
              class="btn btn-primary rounded btn-sm"
              @click="() => {
                command.edit = false
                if (!command.id) commands = commands?.filter((_, i) => i !== commandIndex)
                else if (commands) commands[commandIndex] = commandsBeforeEdit!.find(c => c.id === command.id)!
              }"
            >
              Cancel
            </button>
          </div>
          
          <div>
            <div class="dropdown dropdown-top dropdown-left">
              <label
                tabindex="0"
                class="btn btn-sm rounded btn-error"
              >Delete</label>
              <div
                tabindex="0"
                class="dropdown-content rounded card card-compact w-64 p-2 shadow bg-base-300"
              >
                <h3 class="card-title">
                  Are you sure?
                </h3>
                <div class="card-actions">
                  <div class="btn-group w-full">
                    <button
                      class="btn w-1/2"
                      @click="(e) => {
                        (e.target as HTMLElement).blur()
                      }"
                    >
                      No
                    </button>
                    <button
                      class="btn w-1/2"
                      @click="(e) => {
                        (e.target as HTMLElement).blur()
                        deleteCommand(commandIndex, command.id)
                      }"
                    >
                      Yes
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <button
              class="btn rounded btn-success btn-sm ml-2"
              @click="saveCommand(command, commandIndex)"
            >
              Save
            </button>
          </div>
        </div>

        <div
          v-if="!command.edit"
          class="card-actions justify-end m-3"
        >
          <button
            class="btn roudned btn-primary btn-sm"
            @click="command.edit = true"
          >
            Edit
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
