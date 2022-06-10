<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { UpdateOrCreateCommandDto } from '@tsuwari/api/src/v1/commands/dto/create';
import { CommandPermission, CooldownType } from '@tsuwari/prisma';
import { configure, Form, Field } from 'vee-validate';
import { computed, toRef } from 'vue';
import * as yup from 'yup';

import type { VariablesList } from '@/dashboard/Commands.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const selectedDashboard = useStore(selectedDashboardStore);

type CommandType = UpdateOrCreateCommandDto

configure({
  validateOnBlur: true, // controls if `blur` events should trigger validation with `handleChange` handler
  validateOnChange: true, // controls if `change` events should trigger validation with `handleChange` handler
  validateOnInput: true, // controls if `input` events should trigger validation with `handleChange` handler
  validateOnModelUpdate: true, // controls if `update:modelValue` events should trigger validation with `handleChange` handler
});

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

const props = defineProps<{ 
  command: CommandType,
  commands: CommandType[]
  variablesList: VariablesList
}>();

const command = toRef(props, 'command');
const commands = toRef(props, 'commands');
const emit = defineEmits<{
  (e: 'delete', index: number): void
  (e: 'save', index: number): void
}>();

const schema = computed(() => yup.object({
  name: yup.string().min(1, 'Name cannot be empty')
    .test(
      'unique-name',
      (d) => `Name "${d.value}" already used for other command.`,
      (v) => {
        const otherCommands = commands.value.filter(c => c.id !== command.value.id);

         if (otherCommands?.some(c => c.name === v)) {
          return false;
        }

        if (otherCommands?.some(c => c.aliases?.some(aliase => aliase === v))) {
          return false;
        }

      return true;
      },
    ),
  cooldown: yup.number().notRequired().min(5, 'Cooldown cannot be lower then 5 seconds.'),
  permission: yup.mixed().oneOf(Object.values(perms)),
  aliases: yup.array<Array<string>>().optional().of(
    yup.string().test((v) => {
      if (!v) return false;
      const otherCommands = commands.value.filter(c => c.id !== command.value.id);

      if (otherCommands?.some(c => c.name === v)) {
        return false;
      }

      if (otherCommands?.some(c => c.aliases?.some(aliases => aliases.includes(v)))) {
        return false;
      }

      return true;
    }),
  ),
  commands: yup.array().required('Responses cannot be empty.'),
}));
  
async function deleteCommand() {
  const index = commands.value.indexOf(command.value);
  if (command.value.id) {
    await api.delete(`/v1/channels/${selectedDashboard.value.channelId}/commands/${command.value.id}`);
  }

  emit('delete', index);
}

async function saveCommand() {
  const index = commands.value.indexOf(command.value);
  let data: CommandType;
  if (command.value.id) {
    const request = await api.put(`/v1/channels/${selectedDashboard.value.channelId}/commands/${command.value.id}`, command.value);  
    data = request.data;
  } else {
    const request = await api.post(`/v1/channels/${selectedDashboard.value.channelId}/commands`, command.value);
    data = request.data;
  }

  if (commands.value && commands.value[index]) {
    commands.value[index] = data;
    emit('save', index);
  }
}

function changeCommandResponse(index: number, value: string) {
  command.value.responses[index].text = value;
}

const consoleLog = console.log;
</script>

<template>
  <div class="p-4">
    <Form
      v-slot="{ errors }"
      :validation-schema="schema"
      @submit="saveCommand"
    > 
      <div
        v-for="error of errors"
        :key="error"
        class="bg-red-100 rounded-lg py-5 px-6 mb-4 text-base text-red-700 mb-3"
        role="alert"
      >
        {{ error }}
      </div>
      <div
        class="grid grid-cols-2 gap-1"
      >
        <div>
          <div class="label">
            <span class="label-text">Name</span>
          </div>
          <Field
            v-model="command.name"
            name="name"
            type="text"
            placeholder="uptime"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>


        <div>
          <span class="label text-center">Cooldown</span>
            
          <div class="grid grid-cols-2">
            <Field
              v-model.number="command.cooldown"
              name="cooldown"
              as="input"
              type="number"
              placeholder="0"
              class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-4/5 input-sm"
            />
                
            <Field
              v-model="command.cooldownType"
              name="cooldownType"
              as="select"
              class="form-control px-3 py-1.5 text-gray-700 rounded select select-sm w-full"
            >
              <option
                v-for="type of Object.entries(cooldownType)"
                :key="type[0]"
                :value="type[1]"
              >
                {{ type[0] }}
              </option>
            </Field>
          </div>
        </div>

        <div>
          <span class="label">Permission</span>
          <Field
            v-model="command.permission"
            as="select"
            name="permission"
            class="form-control px-3 py-1.5 text-gray-700 rounded select select-sm w-full"
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
          </Field>
        </div>

        <div>
          <span class="label">Description (optional)</span>
          <Field
            v-model="command.description"
            name="description"
            type="text"
            placeholder="great command ;)"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>

        <div class="col-span-2">
          <span class="label">
            <span>Responses
              <a
                class="px-2 py-0.5 inline-block bg-green-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:shadow-lg focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 cursor-pointer ease-in-out"
                @click="command.responses.push({ text: '' })"
              >
                +
              </a>
            </span>
          </span>

          <!-- max-h-40 scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600-->
          <div class="input-group min-h-[150px] grid grid-cols-1 pt-1 gap-2">
            <div
              v-for="_response, responseIndex in command.responses"
              :key="responseIndex"
              class="flex flex-wrap max-h-max items-stretch mb-4 relative dropdown relative"
              style="width: 99%;"
            >
              <div 
                :id="'dropdownMenuButton' + responseIndex"
                class="form-control dropdown-toggle px-3 py-1.5 text-gray-700 rounded input input-bordered w-[90%] input-sm bg-white rounded-r-none" 
                contenteditable
                data-bs-toggle="dropdown"
                aria-expanded="false"
                @input="(payload) => changeCommandResponse(responseIndex, (payload.target! as HTMLElement).innerText)"
              >
                {{ command.responses[responseIndex].text }}
              </div>
              <ul
                class="
                  dropdown-menu w-[90%] absolute hidden bg-white text-base z-50 float-left py-2 list-none text-left rounded-lg shadow-lg mt-1 hidden m-0 bg-clip-padding border-none bg-gray-800 max-h-52 scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600
                "
                :aria-labelledby="'dropdownMenuButton' + responseIndex"
              >
                <h6
                  class="
                    text-gray-500
                    font-semibold
                    text-sm
                    py-2
                    px-4
                    block
                    w-full
                    whitespace-nowrap
                    bg-transparent
                  "
                >
                  Variables
                </h6>
                <li
                  v-for="variable of variablesList"
                  :key="variable.name"
                >
                  <a
                    class="
                      dropdown-item
                      text-sm
                      py-2
                      px-4
                      font-normal
                      block
                      w-full
                      whitespace-nowrap
                      bg-transparent
                      text-white
                      hover:text-gray-700
                      hover:bg-gray-100
                    "
                    @click="() => {
                      consoleLog(variable)
                      command.responses[responseIndex].text += ` $(${variable.example ? variable.example : variable.name})`;
                    }"
                  >{{ variable.description ?? variable.name }}</a>
                </li>
              </ul>
              <div
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
              <a
                class="px-2 py-0.5 inline-block bg-green-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md  hover:shadow-lg focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out cursor-pointer"
                @click="command.aliases?.push('')"
              >
                +
              </a>
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
                class="flex-shrink rounded-r-none flex-grow flex-auto leading-normal w-px flex-1 border border-grey-light text-gray-700 rounded px-3 py-1.5 relative"
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

      <div class="flex justify-between mt-5">
        <div />
        <div>
          <button
            v-if="command.id"
            type="button"
            class="inline-block px-6 py-2.5 bg-red-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-red-700 hover:shadow-lg focus:bg-red-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-red-800 active:shadow-lg transition duration-150 ease-in-out"
            @click="deleteCommand"
          >
            Delete
          </button>
          <button
            type="submit"
            class="inline-block ml-2 px-6 py-2.5 bg-green-500 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-green-600 hover:shadow-lg focus:bg-green-600 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-green-700 active:shadow-lg transition duration-150 ease-in-out"
          >
            Save
          </button>
        </div>
      </div>
    </Form>
  </div>
</template>

<style scoped>
input, select {
  @apply border-inherit
}
input:disabled, select:disabled {
  @apply bg-zinc-400 opacity-100 border-transparent
}
</style>