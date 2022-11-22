<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { Form, Field } from 'vee-validate';
import { Ref, toRef } from 'vue';
import { useI18n } from 'vue-i18n';

import { CounterType } from '@/dashboard/WordsCounters.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const props =
  defineProps<{
    counter: CounterType;
    counters: CounterType[];
    countersBeforeEdit: CounterType[];
  }>();

const counter = toRef(props, 'counter');
const counters = toRef(props, 'counters');
const countersBeforeEdit = toRef(props, 'countersBeforeEdit');
const selectedDashboard = useStore(selectedDashboardStore);
const { t } = useI18n();
const emit =
  defineEmits<{
    (e: 'delete', index: number): void;
    (e: 'cancelEdit', counter: Ref<CounterType>): void;
  }>();

async function saveCounter() {
  const index = counters.value.indexOf(counter.value);

  let data;

  if (counter.value.id) {
    const request = await api.put(
      `/v1/channels/${selectedDashboard.value.channelId}/words_counters/${counter.value.id}`,
      counter.value,
    );
    data = request.data;
  } else {
    const request = await api.post(
      `/v1/channels/${selectedDashboard.value.channelId}/words_counters`,
      counter.value,
    );
    data = request.data;
  }

  counters.value[index] = data;
}

async function deletecounter() {
  const index = counters.value.indexOf(counter.value);
  if (counter.value.id) {
    await api.delete(
      `/v1/channels/${selectedDashboard.value.channelId}/words_counters/${counter.value.id}`,
    );
  }

  emit('delete', index);
}

function cancelEdit() {
  emit('cancelEdit', counter);
}
</script>

<template>
  <div class="p-4">
    <Form v-slot="{ errors }" @submit="saveCounter">
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
          <p>Status</p>
          <input
            id="commandVisibility"
            v-model="counter.enabled"
            :disabled="!counter.edit"
            class="
              align-top
              appearance-none
              bg-contain bg-gray-300 bg-no-repeat
              cursor-pointer
              float-left
              focus:outline-none
              form-check-input
              h-5
              rounded-full
              shadow
              w-9
            "
            type="checkbox"
            role="switch"
          />
        </div>
      </div>
      <div class="gap-1 grid grid-cols-1">
        <div class="mt-5">
          <div class="label mb-1">
            <span class="label-text">Phrase</span>
          </div>
          <Field
            v-model.lazy.trim="counter.phrase"
            name="phrase"
            as="input"
            type="text"
            :placeholder="t('pages.wordsCounters.card.phrase.placeholder')"
            :disabled="!counter.edit"
            class="
              form-control
              input input-bordered input-sm
              px-3
              py-1.5
              rounded
              text-gray-700
              w-full
            "
          />
        </div>

        <div class="mt-5">
          <div class="label mb-1">
            <span class="label-text">Count</span>
          </div>
          <Field
            v-model.number.lazy="counter.counter"
            name="count"
            as="input"
            type="text"
            :disabled="!counter.edit"
            class="
              form-control
              input input-bordered input-sm
              px-3
              py-1.5
              rounded
              text-gray-700
              w-full
            "
          />
        </div>
      </div>

      <div v-if="counter.id" class="mt-2 p-4 rounded-md bg-cyan-800">
        You can use that counter in commands/keywords/timers via
        <b>{{ `$(words.counter|${counter.id})` }}</b> variable
      </div>

      <div class="flex justify-between mt-5">
        <div>
          <button
            v-if="!counter.edit"
            type="button"
            class="
              bg-purple-600
              duration-150
              ease-in-out
              focus:outline-none
              focus:ring-0
              font-medium
              hover:bg-putple-700
              inline-block
              leading-tight
              px-6
              py-2.5
              rounded
              shadow
              text-xs
              transition
              uppercase
            "
            @click="
              () => {
                counter.edit = true;
                if (counter.id) countersBeforeEdit?.push(JSON.parse(JSON.stringify(counter)));
              }
            "
          >
            {{ t('buttons.edit') }}
          </button>
          <button
            v-else
            type="button"
            class="
              bg-purple-600
              duration-150
              ease-in-out
              focus:outline-none
              focus:ring-0
              font-medium
              hover:bg-purple-700
              inline-block
              leading-tight
              px-6
              py-2.5
              rounded
              shadow
              text-white text-xs
              transition
              uppercase
            "
            @click="cancelEdit"
          >
            {{ t('buttons.cancel') }}
          </button>
        </div>
        <div v-if="counter.edit" class="flex md:flex-none ml-1">
          <button
            v-if="counter.id"
            type="button"
            class="
              bg-red-600
              duration-150
              ease-in-out
              focus:outline-none
              focus:ring-0
              font-medium
              hover:bg-red-700
              inline-block
              leading-tight
              px-6
              py-2.5
              rounded
              shadow
              text-white text-xs
              transition
              uppercase
            "
            @click="deletecounter"
          >
            {{ t('buttons.delete') }}
          </button>
          <button
            type="submit"
            class="
              bg-green-600
              duration-150
              ease-in-out
              focus:outline-none
              focus:ring-0
              font-medium
              hover:bg-green-700
              inline-block
              leading-tight
              ml-1
              px-6
              py-2.5
              rounded
              shadow
              text-white text-xs
              transition
              uppercase
            "
          >
            {{ t('buttons.save') }}
          </button>
        </div>
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
