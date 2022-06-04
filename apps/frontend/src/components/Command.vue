<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { UpdateOrCreateCommandDto } from '@tsuwari/api/src/v1/commands/dto/create';
import { CommandPermission, CooldownType } from '@tsuwari/prisma';
import { configure, Form, Field } from 'vee-validate';
import { computed, toRef } from 'vue';
import * as yup from 'yup';

import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const selectedDashboard = useStore(selectedDashboardStore);

type CommandType = UpdateOrCreateCommandDto & { 
  edit?: boolean
}

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
  commandsBeforeEdit: CommandType[]
}>();

const command = toRef(props, 'command');
const commands = toRef(props, 'commands');
const commandsBeforeEdit = toRef(props, 'commandsBeforeEdit');

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

        if (otherCommands?.some(c => c.aliases?.some(aliases => aliases.includes(v!)))) {
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
  }));
  
async function deleteCommand() {
  const index = commands.value.indexOf(command.value);
  if (command.value.id) {
    await api.delete(`/v1/channels/${selectedDashboard.value.channelId}/commands/${command.value.id}`);
  }

  if (commands.value) {
    commands.value = commands.value.filter((_, i) => i !== index);
  }
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

    const editableCommand = commandsBeforeEdit.value?.find(c => c.id === data.id);
    if (editableCommand) {
      commandsBeforeEdit.value?.splice(commandsBeforeEdit.value.indexOf(editableCommand));
    }
  }
}

function cancelEdit() {
  const index = commands.value.indexOf(command.value);
  if (command.value.id && commands.value) {
    const editableCommand = commandsBeforeEdit.value?.find(c => c.id === command.value.id);
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
            :disabled="!command.edit"
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
              :disabled="!command.edit"
              class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-4/5 input-sm"
            />
                
            <Field
              v-model="command.cooldownType"
              name="cooldownType"
              as="select"
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
            </Field>
          </div>
        </div>

        <div>
          <span class="label">Permission</span>
          <Field
            v-model="command.permission"
            as="select"
            name="permission"
            :disabled="!command.edit"
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
            :disabled="!command.edit"
            placeholder="great command ;)"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>

        <div class="col-span-2">
          <span class="label">
            <span>Responses
              <a
                v-if="command.edit"
                class="px-2 py-0.5 inline-block bg-green-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:shadow-lg focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 cursor-pointer ease-in-out"
                @click="command.responses.push({ text: '' })"
              >
                +
              </a>
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
                v-model="command.responses[responseIndex].text"
                type="text"
                :disabled="!command.edit"
                class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1 border text-gray-700 rounded px-3 py-1.5 relative"
                placeholder="command response"
                :class="{ 'rounded-r-none': command.edit }"
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
              <a
                v-if="command.edit"
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
                :disabled="!command.edit"
                type="text"
                class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1 border border-grey-light text-gray-700 rounded px-3 py-1.5 relative"
                :class="{ 'rounded-r-none': command.edit }"
              >
              <div
                v-if="command.edit"
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
            @click="cancelEdit"
          >
            Cancel
          </button>
        </div>
        <div v-if="command.edit">
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