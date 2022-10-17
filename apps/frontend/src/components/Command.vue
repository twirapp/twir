<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { UpdateOrCreateCommandDto } from '@tsuwari/api/src/v1/commands/dto/create';
import { CommandPermission, CooldownType } from '@tsuwari/typeorm/entities/ChannelCommand';
import { isNumber } from '@vueuse/core';
import { configure, Form, Field } from 'vee-validate';
import { computed, toRef } from 'vue';
import { useI18n } from 'vue-i18n';
import { useToast } from 'vue-toastification';
import * as yup from 'yup';

import Add from '@/assets/buttons/add.svg?component';
import Remove from '@/assets/buttons/remove.svg?component';
import Error from '@/assets/icons/error.svg?component';
import MyBtn from '@/components/elements/MyBtn.vue';
import type { VariablesList } from '@/dashboard/Commands.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const selectedDashboard = useStore(selectedDashboardStore);

type CommandType = UpdateOrCreateCommandDto & { default?: boolean };

configure({
  validateOnBlur: true, // controls if `blur` events should trigger validation with `handleChange` handler
  validateOnChange: true, // controls if `change` events should trigger validation with `handleChange` handler
  validateOnInput: true, // controls if `input` events should trigger validation with `handleChange` handler
  validateOnModelUpdate: true, // controls if `update:modelValue` events should trigger validation with `handleChange` handler
});

const perms = {
  Broadcaster: 'BROADCASTER',
  // eslint-disable-next-line quotes
  "Moderator's": 'MODERATOR',
  // eslint-disable-next-line quotes
  "Vip's": 'VIP',
  // eslint-disable-next-line quotes
  "Subscriber's": 'SUBSCRIBER',
  Followers: 'FOLLOWER',
  Viewers: 'VIEWER',
} as { [x: string]: CommandPermission };

const cooldownType = {
  Global: 'GLOBAL',
  'Per user': 'PER_USER',
} as { [x: string]: CooldownType };

const props = defineProps<{
  command: CommandType;
  commands: CommandType[];
  variablesList: VariablesList;
}>();

const command = toRef(props, 'command');
const commands = toRef(props, 'commands');
const emit = defineEmits<{
  (e: 'delete', index: number): void;
  (e: 'save', index: number): void;
}>();
const { t } = useI18n({
  useScope: 'global',
});
const toast = useToast();

const schema = computed(() =>
  yup.object({
    name: yup
      .string()
      .min(1, 'Name cannot be empty')
      .test(
        'unique-name',
        (d) => `Name "${d.value}" already used for other command.`,
        (v) => {
          const otherCommands = commands.value.filter((c) => c.id !== command.value.id);

          if (otherCommands?.some((c) => c.name === v)) {
            return false;
          }

          if (otherCommands?.some((c) => c.aliases?.some((aliase) => aliase === v))) {
            return false;
          }

          return true;
        },
      ),
    cooldown: yup
      .number()
      .notRequired()
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
    aliases: yup
      .array<Array<string>>()
      .optional()
      .of(
        yup.string().test((v) => {
          if (!v) return false;
          const otherCommands = commands.value.filter((c) => c.id !== command.value.id);

          if (otherCommands?.some((c) => c.name === v)) {
            return false;
          }

          if (otherCommands?.some((c) => c.aliases?.some((aliases) => aliases.includes(v)))) {
            return false;
          }

          return true;
        }),
      ),
  }),
);

async function deleteCommand() {
  const index = commands.value.indexOf(command.value);
  if (command.value.id) {
    await api.delete(
      `/v1/channels/${selectedDashboard.value.channelId}/commands/${command.value.id}`,
    );
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
      {
        ...command.value,
        responses: command.value.responses.filter((r) => r.text),
      },
    );
    data = request.data;
    toast.success('Command updated');
  } else {
    const request = await api.post(`/v1/channels/${selectedDashboard.value.channelId}/commands`, {
      ...command.value,
      responses: command.value.responses.filter((r) => r.text),
    });
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
  <div v-if="command" class="p-4">
    <Form v-slot="{ errors }" :validation-schema="schema" @submit="saveCommand">
      <div
        v-for="error of errors"
        :key="error"
        class="bg-red-600 flex mb-4 px-6 py-2 rounded text-white"
        role="alert">
        <Error />
        <p>{{ error }}</p>
      </div>
      <div class="flex justify-end">
        <div class="flex form-switch space-x-2">
          <p>{{ t('pages.commands.card.status.title') }}</p>
          <input
            id="commandVisibility"
            v-model="command.enabled"
            class="align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
            type="checkbox"
            role="switch" />
        </div>
      </div>

      <div>
        <div>
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.commands.card.name.title') }}</span>
          </div>
          <Field
            v-model.trim="command.name"
            name="name"
            type="text"
            :placeholder="t('pages.commands.card.name.placeholder')"
            class="form-control input input-bordered input-sm px-3 py-1.5 rounded text-gray-700 w-full" />
        </div>
        <div class="mt-5">
          <span class="label text-center">{{ t('pages.commands.card.cooldown.title') }}</span>

          <div class="grid grid-cols-2 mt-1">
            <Field
              v-model.number="command.cooldown"
              name="cooldown"
              as="input"
              type="number"
              placeholder="0"
              class="form-control input input-bordered input-sm px-3 py-1.5 rounded text-gray-700 w-4/5" />

            <Field
              v-model.trim="command.cooldownType"
              name="cooldownType"
              as="select"
              class="form-control px-3 py-1.5 rounded select select-sm text-gray-700 w-full">
              <option v-for="type of Object.entries(cooldownType)" :key="type[0]" :value="type[1]">
                {{ t(`pages.commands.card.cooldown.type.${type[1]}`) }}
              </option>
            </Field>
          </div>
        </div>
      </div>

      <div class="gap-1 grid grid-cols-1 md:grid-cols-2">
        <div class="col-span-2 mt-5">
          <span class="flex items-center label">
            <span>{{ t('pages.commands.card.responses.title') }} </span>
            <MyBtn
              :show="!command.default"
              color="green"
              size="small"
              class="ml-1"
              @click="command.responses.push({ text: '', order: command.responses.length })"
              ><Add
            /></MyBtn>
          </span>

          <div v-if="!command.default" class="gap-1 grid grid-cols-1 input-group mшx-h-[5px] pt-1">
            <div
              v-for="(_response, responseIndex) in command.responses"
              :key="responseIndex"
              class="dropdown flex items-stretch max-h-max mb-1 relative">
              <div class="flex items-stretch relative w-full">
                <div
                  :id="'dropdownMenuButton' + responseIndex"
                  class="bg-white dropdown-toggle form-control input input-bordered input-sm px-3 py-1.5 rounded rounded-r-none text-gray-700 w-full"
                  contenteditable
                  data-bs-toggle="dropdown"
                  aria-expanded="false"
                  @input="(payload) => changeCommandResponse(responseIndex, (payload.target! as HTMLElement).innerText.trim())">
                  {{ command.responses[responseIndex].text }}
                </div>

                <span
                  class="bg-red-600 border-0 border-grey-light border-l-0 flex hover:bg-red-700 items-center leading-normal px-4 py-1.5 rounded rounded-l-none text-grey-dark text-sm whitespace-no-wrap"
                  @click="command.responses?.splice(responseIndex, 1)"
                  ><Remove
                /></span>

                <ul
                  class="absolute bg-[#393636] bg-clip-padding border-none dropdown-menu float-left hidden list-none m-0 max-h-52 mt-1 overflow-auto py-2 rounded scrollbar shadow text-base text-left text-white w-[90%] z-50"
                  :aria-labelledby="'dropdownMenuButton' + responseIndex">
                  <h6
                    class="bg-transparent block font-semibold px-4 py-2 text-sm w-full whitespace-nowrap">
                    Variables
                  </h6>
                  <li v-for="variable of variablesList" :key="variable.name">
                    <a
                      class="bg-transparent block dropdown-item font-normal hover:bg-[#4f4a4a] px-4 py-2 text-sm text-white w-full whitespace-nowrap"
                      @click="
                        () => {
                          command.responses[responseIndex].text += ` $(${
                            variable.example ? variable.example : variable.name
                          })`;
                        }
                      "
                      >{{ variable.description ?? variable.name }}</a
                    >
                  </li>
                </ul>
              </div>
            </div>
          </div>

          <div v-else class="bg-[#ED4245] flex px-6 py-2 rounded text-white" role="alert">
            {{ t('pages.commands.card.responses.builtInAlert') }}
          </div>
        </div>

        <div class="mt-5">
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.commands.card.permission.title') }}</span>
          </div>
          <Field
            v-model.trim="command.permission"
            as="select"
            name="permission"
            class="form-control px-3 py-1.5 rounded select select-sm text-gray-700">
            <option disabled selected>
              {{ t('pages.commands.card.permission.selectPlaceholder') }}
            </option>
            <option
              v-for="permission of Object.entries(perms)"
              :key="permission[0]"
              :value="permission[1]">
              {{ permission[0] }}
            </option>
          </Field>
        </div>
      </div>

      <div class="mt-5">
        <div class="flex form-check justify-between">
          <label class="form-check-label inline-block" for="keepOrder">{{
            t('pages.commands.card.keepOrder.title')
          }}</label>

          <div class="form-switch">
            <input
              id="keepOrder"
              v-model="command.keepOrder"
              class="align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
              type="checkbox"
              role="switch" />
          </div>
        </div>
      </div>

      <div class="mt-5">
        <div class="flex form-check justify-between">
          <label class="form-check-label inline-block" for="commandVisibility">{{
            t('pages.commands.card.visible.title')
          }}</label>

          <div class="form-switch">
            <input
              id="commandVisibility"
              v-model="command.visible"
              class="align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
              type="checkbox"
              role="switch" />
          </div>
        </div>
      </div>

      <div class="mt-5">
        <div class="flex form-check justify-between">
          <label class="form-check-label inline-block" for="commandIsReply">{{
            t('pages.commands.card.isReply.title')
          }}</label>

          <div class="form-switch">
            <input
              id="commandIsReply"
              v-model="command.isReply"
              class="align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
              type="checkbox"
              role="switch" />
          </div>
        </div>
      </div>

      <div class="mt-5">
        <div class="label mb-1">
          <span class="label-text">{{ t('pages.commands.card.description.title') }}</span>
        </div>
        <Field
          v-model.trim="command.description"
          name="description"
          as="textarea"
          :placeholder="t('pages.commands.card.description.placeholder')"
          class="bg-clip-padding bg-white block border border-gray-300 border-solid ease-in-out focus:bg-white focus:border-blue-600 focus:outline-none focus:text-gray-700 font-normal form-control m-0 px-3 py-1.5 rounded text-base text-gray-700 transition w-full"
          rows="2" />
      </div>

      <div class="col-span-2 mt-5">
        <span class="flex items-center label">
          <span>{{ t('pages.commands.card.aliases.title') }}</span>
          <MyBtn color="green" size="small" class="ml-1" @click="command.aliases?.push('')"
            ><Add
          /></MyBtn>
        </span>

        <div
          class="gap-2 grid grid-cols-1 input-group lg:grid-cols-2 max-h-40 md:grid-cols-2 overflow-auto pr-2 pt-1 scrollbar sm:grid-cols-2 xl:grid-cols-3">
          <div
            v-for="(aliase, aliaseIndex) in command.aliases"
            :key="aliase"
            class="flex flex-wrap items-stretch relative">
            <input
              v-model.lazy.trim="command.aliases![aliaseIndex]"
              type="text"
              class="border border-grey-light flex-grow flex-shrink leading-normal px-3 py-1.5 relative rounded rounded-r-none text-gray-700 w-px" />
            <div class="cursor-pointer flex" @click="command.aliases?.splice(aliaseIndex, 1)">
              <span
                class="bg-red-600 border-0 border-grey-light border-l-0 flex hover:bg-red-700 items-center leading-normal px-5 py-1.5 rounded rounded-l-none text-grey-dark text-sm whitespace-no-wrap"
                ><Remove
              /></span>
            </div>
          </div>
        </div>
      </div>

      <div
        class="flex flex-col justify-end md:flex-row md:space-x-2 md:space-y-0 mt-5 space-y-2 w-full">
        <MyBtn :show="Boolean(command.id && !command.default)" color="red" @click="deleteCommand">
          {{ t('buttons.delete') }}
        </MyBtn>

        <MyBtn color="green" type="submit">
          {{ t('buttons.save') }}
        </MyBtn>
      </div>
    </Form>
  </div>
</template>

<style scoped>
input,
select {
  @apply border-inherit;
}
input:disabled,
select:disabled {
  @apply bg-zinc-400 opacity-100 border-transparent;
}
</style>
