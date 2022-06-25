<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { Form, Field } from 'vee-validate';
import { computed, toRef } from 'vue';
import { useI18n } from 'vue-i18n';
import * as yup from 'yup';

import { GreeTingType } from '@/dashboard/Greetings.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const props = defineProps<{
  greeting: GreeTingType,
  greetings: GreeTingType[],
  greetingsBeforeEdit: GreeTingType[]
}>();

const greeting = toRef(props, 'greeting');
const greetings = toRef(props, 'greetings');
const greetingsBeforeEdit = toRef(props, 'greetingsBeforeEdit');
const selectedDashboard = useStore(selectedDashboardStore);
const { t } = useI18n({
  useScope: 'global',
});

const emit = defineEmits<{
  (e: 'delete', index: number): void
}>();

const schema = computed(() => yup.object({
  username: yup.string().required(),
  text: yup.string().required(),
}));

async function saveGreeting() {
  const index = greetings.value.indexOf(greeting.value);

  let data;
  if (props.greeting.id) {
    const request = await api.put(`/v1/channels/${selectedDashboard.value.channelId}/greetings/${greeting.value.id}`, greeting.value);  
    data = request.data;
  } else {
    const request = await api.post(`/v1/channels/${selectedDashboard.value.channelId}/greetings`, greeting.value);
    data = request.data;
  }

  greetings.value[index] = data;
}

async function deleteGreeting() {
  const index = greetings.value.indexOf(greeting.value);
  if (greeting.value.id) {
    await api.delete(`/v1/channels/${selectedDashboard.value.channelId}/greetings/${greeting.value.id}`);
  }

  emit('delete', index);
}

function cancelEdit() {
  const index = greetings.value.indexOf(greeting.value);
  if (greeting.value.id && greetings.value) {
    const editableCommand = greetingsBeforeEdit.value?.find(c => c.id === greeting.value.id);
    if (editableCommand) {
      greetings.value[index] = {
        ...editableCommand,
        edit: false,
      };
      greetingsBeforeEdit.value?.splice(greetingsBeforeEdit.value.indexOf(editableCommand), 1);
    }
  } else {
    greetings.value?.splice(index, 1);
  }
}
</script>

<template>
  <div class="p-4">
    <Form
      v-slot="{ errors }"
      :validation-schema="schema"
      @submit="saveGreeting"
    >
      <div
        v-for="error of errors"
        :key="error"
        class="bg-red-600 rounded py-5 px-6 mb-4 text-white"
        role="alert"
      >
        {{ error }}
      </div>
      <div
        class="grid grid-cols-1 gap-1"
      >
        <div>
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.greetings.username.title') }}</span>
          </div>
          <Field
            v-model="greeting.username"
            name="username"
            as="input" 
            type="text"
            :placeholder="t('pages.greetings.username.placeholder')"
            :disabled="!greeting.edit"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>

        <div class="mt-5">
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.greetings.message.title') }}</span>
          </div>
          <Field
            v-model="greeting.text"
            name="text"
            as="input" 
            type="text"
            :placeholder="t('pages.greetings.message.placeholder')"
            :disabled="!greeting.edit"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>
      </div>

      <div class="flex justify-between mt-5">
        <div>
          <button
            v-if="!greeting.edit"
            type="button"
            class="inline-block px-6 py-2.5 bg-purple-600 font-medium text-xs leading-tight uppercase rounded shadow hover:bg-purple-700 focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
            @click="() => {
              greeting.edit = true;
              if (greeting.id) greetingsBeforeEdit?.push(JSON.parse(JSON.stringify(greeting)))
            }"
          >
            {{ t('buttons.edit') }}
          </button>
          <button
            v-else
            class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-purple-700  focus:outline-none focus:ring-0  transition duration-150 ease-in-out"
            @click="cancelEdit"
          >
            {{ t('buttons.cancel') }}
          </button>
        </div>
        <div v-if="greeting.edit">
          <button
            v-if="greeting.id"
            type="button"
            class="inline-block px-6 py-2.5 bg-red-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-red-700  focus:outline-none focus:ring-0   transition duration-150 ease-in-out"
            @click="deleteGreeting"
          >
            {{ t('buttons.delete') }}
          </button>
          <button
            type="submit"
            class="inline-block ml-2 px-6 py-2.5 bg-green-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-green-700  focus:outline-none focus:ring-0   transition duration-150 ease-in-out"
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