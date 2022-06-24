<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { UpdateOrCreateCommandDto } from '@tsuwari/api/src/v1/commands/dto/create';
import { CommandPermission, CooldownType } from '@tsuwari/prisma';
import { isNumber } from '@vueuse/core';
import axios from 'axios';
import { configure, Form, Field } from 'vee-validate';
import { computed, toRef } from 'vue';
import { useI18n } from 'vue-i18n';
import { useToast } from 'vue-toastification';
import * as yup from 'yup';

import Add from '@/assets/buttons/add.svg';
import Remove from '@/assets/buttons/remove.svg';
import type { VariablesList } from '@/dashboard/Commands.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const selectedDashboard = useStore(selectedDashboardStore);

type CommandType = UpdateOrCreateCommandDto & { default?: boolean }

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
const { t } = useI18n({
  useScope: 'global',
});
const toast = useToast();

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
  cooldown: yup.number().notRequired()
    .test(
      'cooldown', 
      () => `Cooldown cannot be lower then 5 seconds.`,
      (v) => {
        if (typeof v === 'undefined' || !isNumber(v)) return false;
        if (command.value.default) return true;

        return v >= 5;
      },
  ),
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
    toast.success('Command deleted');
  }

  emit('delete', index);
}

async function saveCommand() {
  const index = commands.value.indexOf(command.value);
  let data: CommandType;

  if (command.value.id) {
    const request = await api.put(
      `/v1/channels/${selectedDashboard.value.channelId}/commands/${command.value.id}`,
      command.value,
    );  
    data = request.data;
    toast.success('Command updated');
  } else {
    const request = await api.post(
      `/v1/channels/${selectedDashboard.value.channelId}/commands`, 
      command.value,
    );
    data = request.data;
    toast.success('Command created');
  }

  if (commands.value && commands.value[index]) {
    commands.value[index] = data;
    emit('save', index);
  }
}

function changeCommandResponse(index: number, value: string) {
  command.value.responses[index].text = value;
}
</script>

<template>
  <div
    v-if="command"
    class="p-4"
  >
    <Form
      v-slot="{ errors }"
      :validation-schema="schema"
      @submit="saveCommand"
    > 
      <div
        v-for="error of errors"
        :key="error"
        class="bg-red-600 rounded py-2 px-6 mb-4 text-white flex"
        role="alert"
      >
        <svg
          class="w-6 h-6 mr-2"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
        /></svg>
        <p>{{ error }}</p>
      </div>
      <div
        class="grid grid-cols-2 gap-1"
      >
        <div>
          <div class="label">
            <span class="label-text">{{ t('pages.commands.card.name.title') }}</span>
          </div>
          <Field
            v-model="command.name"
            name="name"
            type="text"
            :placeholder="t('pages.commands.card.name.placeholder')"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>


        <div>
          <span class="label text-center">{{ t('pages.commands.card.cooldown.title') }}</span>
            
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
                {{ t(`pages.commands.card.cooldown.type.${type[1]}`) }}
              </option>
            </Field>
          </div>
        </div>

        <div>
          <span class="label">{{ t('pages.commands.card.permission.title') }}</span>
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
              {{ t('pages.commands.card.permission.selectPlaceholder') }}
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
          <span class="label">{{ t('pages.commands.card.description.title') }}</span>
          <Field
            v-model="command.description"
            name="description"
            type="text"
            :placeholder="t('pages.commands.card.description.placeholder')"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>

        <div class="col-span-2">
          <span class="label flex items-center">
            <span>{{ t('pages.commands.card.responses.title') }}
            </span>
            <span
              v-if="!command.default"
              class="px-1 ml-1 py-1 inline-block bg-green-600 hover:bg-green-700 text-white font-medium text-xs leading-tight uppercase rounded shadow   focus:outline-none focus:ring-0 ansition duration-150 cursor-pointer ease-in-out"
              @click="command.responses.push({ text: '' })"
            >
              <Add />
            </span>
          </span>

          <!-- max-h-40 scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600-->
          <div
            v-if="!command.default"
            class="input-group min-h-[150px] grid grid-cols-1 pt-1 gap-1"
          >
            <div
              v-for="_response, responseIndex in command.responses"
              :key="responseIndex"
              class="flex max-h-max items-stretch mb-1 relative dropdown"
            >
              <div class="flex w-full items-stretch relative">
                <div 
                  :id="'dropdownMenuButton' + responseIndex"
                  class="form-control w-full dropdown-toggle px-3 py-1.5 text-gray-700 rounded input input-bordered input-sm bg-white rounded-r-none" 
                  contenteditable
                  data-bs-toggle="dropdown"
                  aria-expanded="false"
                  @input="(payload) => changeCommandResponse(responseIndex, (payload.target! as HTMLElement).innerText)"
                >
                  {{ command.responses[responseIndex].text }}
                </div>
 
                <span
                  class="flex items-center leading-normal bg-red-600 hover:bg-red-700 rounded rounded-l-none border-0 border-l-0 border-grey-light px-4 py-1.5 whitespace-no-wrap text-grey-dark text-sm"
                  @click="command.responses?.splice(responseIndex, 1)"
                ><Remove /></span>

                <ul
                  class="
                  dropdown-menu w-[90%] absolute bg-white text-base z-50 float-left py-2 list-none text-left rounded shadow mt-1 hidden m-0 bg-clip-padding border-none bg-gray-800 max-h-52 scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600
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
                      hover:bg-gray-200
                    "
                      @click="() => {
                        command.responses[responseIndex].text += ` $(${variable.example ? variable.example : variable.name})`;
                      }"
                    >{{ variable.description ?? variable.name }}</a>
                  </li>
                </ul>
              </div>
            </div>
          </div>
          
          <div
            v-else
            class="bg-[#ED4245] rounded py-2 px-6 mb-4 text-white flex"
            role="alert"
          >
            {{ t('pages.commands.card.responses.builtInAlert') }}
          </div>
        </div>

        <div class="col-span-2">
          <span class="label flex items-center">  
            <span>{{ t('pages.commands.card.aliases.title') }}</span>
            <span
              class="items-center ml-1 px-1 py-1 inline-block bg-green-600 hover:bg-green-700 text-white font-medium text-xs leading-tight uppercase rounded shadow    focus:outline-none focus:ring-0 transition duration-150 ease-in-out cursor-pointer"
              @click="command.aliases?.push('')"
            >
              <Add />
            </span>
           
          </span>

          <div class="input-group pt-1 pr-2 grid lg:grid-cols-2 md:grid-cols-2 sm:grid-cols-2 grid-cols-1 xl:grid-cols-3 gap-2 max-h-40 scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600">
            <div
              v-for="aliase, aliaseIndex in command.aliases"
              :key="aliase"
              class="flex flex-wrap items-stretch relative"
            >
              <input
                v-model.lazy="command.aliases![aliaseIndex]"
                type="text"
                class="flex-shrink rounded-r-none flex-grow leading-normal w-px border border-grey-light text-gray-700 rounded px-3 py-1.5 relative"
              >
              <div
                class="flex cursor-pointer"
                @click="command.aliases?.splice(aliaseIndex, 1)"
              >
                <span class="flex items-center leading-normal bg-red-600 hover:bg-red-700 rounded rounded-l-none border-0 border-l-0 border-grey-light px-5 py-1.5 whitespace-no-wrap text-grey-dark text-sm"><Remove /></span>
              </div>
            </div>
          </div>
        </div>
      </div>


      <div class="mt-5 space-y-2 justify-end flex flex-col w-full md:flex-row md:space-x-2 md:space-y-0">
        <button
          v-if="command.id && !command.default"
          class="inline-block px-6 py-2.5 bg-red-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-red-700  focus:outline-none focus:ring-0   transition duration-150 ease-in-out"
          @click="deleteCommand"
        >
          {{ t('buttons.delete') }}
        </button>

        <button
          type="submit"
          class="inline-block px-6 py-2.5 bg-green-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-green-700  focus:outline-none focus:ring-0   transition duration-150 ease-in-out"
        >
          {{ t('buttons.save') }}
        </button>
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