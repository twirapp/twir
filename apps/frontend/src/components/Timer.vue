<script lang="ts" setup>
import MyBtn from '@elements/MyBtn.vue';
import { useStore } from '@nanostores/vue';
import { Timer } from '@tsuwari/prisma';
import type { SetOptional } from 'type-fest';
import { Form, Field } from 'vee-validate';
import { Ref, toRef } from 'vue';
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
  (e: 'delete', index: number): void,
  (e: 'cancelEdit', timer: Ref<TimerType>): void
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
  emit('cancelEdit', timer);
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
        class="bg-red-600 mb-4 px-6 py-5 rounded text-base text-red-700"
        role="alert"
      >
        {{ error }}
      </div>
      <div class="flex justify-end">
        <div class="flex form-switch space-x-2">
          <p>{{ t('pages.timers.card.status.title') }}</p>
          <input
            id="commandVisibility"
            v-model="timer.enabled"
            :disabled="!timer.edit"
            class="align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
            type="checkbox"
            role="switch"
          >
        </div>
      </div>
      <div>
        <div>
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.timers.card.name.title') }}</span>
          </div>
          <Field
            v-model.trim="timer.name"
            name="name"
            as="input" 
            type="text"
            :placeholder="t('pages.timers.card.name.placeholder')"
            :disabled="!timer.edit"
            class="form-control input input-bordered input-sm px-3 py-1.5 rounded text-gray-700 w-full"
          />
        </div>

        <div>
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.timers.card.secondsInterval') }}</span>
          </div>
          <div class="flex justify-between w-[100%]">
            <div class="flex flex-col my-2 w-[75%]">
              <Field
                v-model.number="timer.timeInterval"
                as="input" 
                type="range"
                name="timeInterval"
                :disabled="!timer.edit"
                min="1"
                max="120"
                class="appearance-none bg-transparent focus:outline-none focus:ring-0 focus:shadow-none form-range h-6 p-0 w-full"
              />
              <ul class="flex justify-between px-[10px] w-full">
                <li class="flex justify-center relative">
                  <span class="absolute">1</span>
                </li>
  
                <li class="flex justify-center relative">
                  <span class="absolute">120</span>
                </li>
              </ul>
            </div>
            <div class="w-[18%]">
              <input 
                v-model.number="timer.timeInterval"
                :disabled="!timer.edit"
                class="form-control input input-bordered input-sm px-3 py-1.5 rounded text-gray-700 w-full"
                type="number"
              >
            </div>
          </div>
        </div>

        <div class="mt-5">
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.timers.card.messagesInterval') }}</span>
          </div>
          <Field
            v-model.number="timer.messageInterval"
            name="messagesInterval"
            as="input" 
            type="number"
            :disabled="!timer.edit"
            class="form-control input input-bordered input-sm px-3 py-1.5 rounded text-gray-700 w-full"
          />
        </div>

        <div class="col-span-2 mt-1">
          <span class="flex items-center label">
            <span>{{ t('pages.timers.card.responses') }}</span>
            <span
              v-if="timer.edit"
              class="bg-green-600 cursor-pointer duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-green-700 inline-block leading-tight ml-2 px-1 py-1 rounded shadow text-white text-xs transition uppercase"
              @click="timer.responses.push('')"
            >
              <Add />
          
            </span>
          </span>

          <div class="gap-2 grid grid-cols-1 input-group max-h-40 overflow-auto pt-1 scrollbar">
            <div
              v-for="_response, responseIndex in timer.responses"
              :key="responseIndex"
              class="flex flex-wrap items-stretch relative"
              style="width: 99%;"
            >
              <input
                v-model.trim="timer.responses[responseIndex]"
                type="text"
                :disabled="!timer.edit"
                class="border flex-1 flex-auto flex-grow flex-shrink leading-normal px-3 py-1.5 relative rounded text-gray-700 w-px"
                :placeholder="t('pages.timers.card.response')"
                :class="{ 'rounded-r-none': timer.edit }"
              >
              <div
                v-if="timer.edit"
                class="cursor-pointer flex"
                @click="timer.responses?.splice(responseIndex, 1)"
              >
                <span class="bg-red-600 border-0 border-grey-light border-l-0 flex hover:bg-red-700 items-center leading-normal px-5 py-1.5 rounded rounded-l-none text-grey-dark text-sm whitespace-no-wrap"><Remove /></span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="flex flex-col lg:flex-row lg:justify-between lg:space-x-1 lg:space-y-0 mt-5 space-y-1">
        <div>
          <MyBtn
            v-if="!timer.edit"
            class="w-full"
            color="purple"
            @click="() => {
              timer.edit = true;
              if (timer.id) timersBeforeEdit?.push(JSON.parse(JSON.stringify(timer)))
            }"
          >
            {{ t('buttons.edit') }}
          </MyBtn>

          <MyBtn
            v-else
            class="lg:w-auto w-full"
            color="purple"
            @click="cancelEdit"
          >
            {{ t('buttons.cancel') }}
          </MyBtn>
        </div>
        <div
          v-if="timer.edit"
          class="flex md:flex-none space-x-1"
        >
          <MyBtn
            v-if="timer.id"
            class="lg:w-auto w-1/2"
            color="red"
            @click="deleteTimer"
          >
            {{ t('buttons.delete') }}
          </MyBtn>
          <MyBtn
            class="lg:w-auto w-1/2"
            color="green"
            type="submit"
          >
            {{ t('buttons.save') }}
          </MyBtn>
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