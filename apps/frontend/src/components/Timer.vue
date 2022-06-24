<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { Timer } from '@tsuwari/prisma';
import type { SetOptional } from 'type-fest';
import { Form, Field } from 'vee-validate';
import { toRef } from 'vue';
import { useI18n } from 'vue-i18n';

import Add from '@/assets/buttons/add.svg';
import Remove from '@/assets/buttons/remove.svg';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const selectedDashboard = useStore(selectedDashboardStore);
type TimerType = SetOptional<Omit<Timer, 'responses' | 'channelId'> & { edit?: boolean, responses: string[] }, 'id'>

const props = defineProps<{ 
  timer: TimerType,
  timers: TimerType[]
  timersBeforeEdit: TimerType[]
}>();

const timer = toRef(props, 'timer');
const timers = toRef(props, 'timers');
const timersBeforeEdit = toRef(props, 'timersBeforeEdit');
const { t } = useI18n({
  useScope: 'global',
});

const emit = defineEmits<{
  (e: 'delete', index: number): void
}>();

async function deleteTimer() {
  const index = timers.value.indexOf(timer.value);
  if (timer.value.id) {
    await api.delete(`/v1/channels/${selectedDashboard.value.channelId}/timers/${timer.value.id}`);
  }

  emit('delete', index);
}

async function saveTimer() {
  const index = timers.value.indexOf(timer.value);
  let data: TimerType;
  if (timer.value.id) {
    const request = await api.put(`/v1/channels/${selectedDashboard.value.channelId}/timers/${timer.value.id}`, timer.value);  
    data = request.data;
  } else {
    const request = await api.post(`/v1/channels/${selectedDashboard.value.channelId}/timers`, timer.value);
    data = request.data;
  }

  if (timers.value && timers.value[index]) {
    timers.value[index] = data;

    const editableTimer = timersBeforeEdit.value?.find(c => c.id === data.id);
    if (editableTimer) {
      timersBeforeEdit.value?.splice(timersBeforeEdit.value.indexOf(editableTimer));
    }
  }
}

function cancelEdit() {
  const index = timers.value.indexOf(timer.value);
  if (timer.value.id && timers.value) {
    const editableCommand = timersBeforeEdit.value?.find(c => c.id === timer.value.id);
    if (editableCommand) {
      timers.value[index] = {
        ...editableCommand,
        edit: false,
      };
      timersBeforeEdit.value?.splice(timersBeforeEdit.value.indexOf(editableCommand), 1);
    }
  } else {
    timers.value?.splice(index, 1);
  }
}
</script>

<template>
  <div class="p-2">
    <Form
      v-slot="{ errors }"
      @submit="saveTimer"
    > 
      <div
        v-for="error of errors"
        :key="error"
        class="bg-red-600 rounded py-5 px-6 mb-4 text-base text-red-700"
        role="alert"
      >
        {{ error }}
      </div>
      <div
        class="md:grid grid-cols-2 gap-2"
      >
        <div>
          <div class="label mb-3">
            <span class="label-text">{{ t('pages.timers.card.name.title') }}</span>
          </div>
          <Field
            v-model="timer.name"
            name="name"
            as="input" 
            type="text"
            :placeholder="t('pages.timers.card.name.placeholder')"
            :disabled="!timer.edit"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>

        <div>
          <div class="label mb-3">
            <span class="label-text">{{ t('pages.timers.card.secondsInterval') }}</span>
          </div>
          <Field
            v-model.number="timer.timeInterval"
            as="input" 
            type="number"
            name="timeInterval"
            :disabled="!timer.edit"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>

        <div>
          <div class="label mb-3">
            <span class="label-text">{{ t('pages.timers.card.messagesInterval') }}</span>
          </div>
          <Field
            v-model.number="timer.messageInterval"
            name="messagesInterval"
            as="input" 
            type="number"
            :disabled="!timer.edit"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>

        <div class="col-span-2 mt-1">
          <span class="label flex items-center">
            <span>{{ t('pages.timers.card.responses') }}</span>
            <span
              v-if="timer.edit"
              class="ml-2 px-1 py-1 inline-block bg-green-600 hover:bg-green-700 text-white font-medium text-xs leading-tight uppercase rounded shadow focus:outline-none focus:ring-0  transition duration-150 cursor-pointer ease-in-out"
              @click="timer.responses.push('')"
            >
              <Add />
          
            </span>
          </span>

          <div class="input-group grid grid-cols-1 pt-1 gap-2 max-h-40 scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600">
            <div
              v-for="_response, responseIndex in timer.responses"
              :key="responseIndex"
              class="flex flex-wrap items-stretch relative"
              style="width: 99%;"
            >
              <input
                v-model="timer.responses[responseIndex]"
                type="text"
                :disabled="!timer.edit"
                class="flex-shrink flex-grow flex-auto leading-normal w-px flex-1 border text-gray-700 rounded px-3 py-1.5 relative"
                placeholder="Timer response"
                :class="{ 'rounded-r-none': timer.edit }"
              >
              <div
                v-if="timer.edit"
                class="flex cursor-pointer"
                @click="timer.responses?.splice(responseIndex, 1)"
              >
                <span class="flex items-center leading-normal bg-red-600 hover:bg-red-700 rounded rounded-l-none border-0 border-l-0 border-grey-light px-5 py-1.5 whitespace-no-wrap text-grey-dark text-sm"><Remove /></span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="flex justify-between mt-5">
        <div>
          <button
            v-if="!timer.edit"
            type="button"
            class="inline-block px-6 py-2.5 bg-gray-200 text-gray-700 font-medium text-xs leading-tight uppercase rounded shadow hover:bg-gray-300    focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
            @click="() => {
              timer.edit = true;
              if (timer.id) timersBeforeEdit?.push(JSON.parse(JSON.stringify(timer)))
            }"
          >
            {{ t('buttons.edit') }}
          </button>
          <button
            v-else
            class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-purple-700    focus:outline-none focus:ring-0  transition duration-150 ease-in-out"
            @click="cancelEdit"
          >
            {{ t('buttons.cancel') }}
          </button>
        </div>
        <div v-if="timer.edit">
          <button
            v-if="timer.id"
            type="button"
            class="inline-block px-6 py-2.5 bg-red-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-red-700 focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
            @click="deleteTimer"
          >
            {{ t('buttons.delete') }}
          </button>
          <button
            type="submit"
            class="inline-block ml-2 px-6 py-2.5 bg-green-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-green-700 focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
          >
            {{ t('buttons.save') }}
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