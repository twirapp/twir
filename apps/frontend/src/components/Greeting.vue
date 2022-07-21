<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { Form, Field } from 'vee-validate';
import { computed, Ref, toRef } from 'vue';
import { useI18n } from 'vue-i18n';
import * as yup from 'yup';

import MyBtn from '@/components/elements/MyBtn.vue';
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
  (e: 'delete', index: number): void,
  (e: 'cancelEdit', greeting: Ref<GreeTingType>): void
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
  emit('cancelEdit', greeting);
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
        class="bg-red-600 mb-4 px-6 py-5 rounded text-white"
        role="alert"
      >
        {{ error }}
      </div>
      <div class="flex justify-end">
        <div class="flex form-switch space-x-2">
          <p>{{ t('pages.greetings.card.status.title') }}</p>
          <input
            id="commandVisibility"
            v-model="greeting.enabled"
            :disabled="!greeting.edit"
            class="align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
            type="checkbox"
            role="switch"
          >
        </div>
      </div>
      <div
        class="gap-1 grid grid-cols-1"
      >
        <div>
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.greetings.username.title') }}</span>
          </div>
          <Field
            v-model.trim="greeting.username"
            name="username"
            as="input" 
            type="text"
            :placeholder="t('pages.greetings.username.placeholder')"
            :disabled="!greeting.edit"
            class="form-control input input-bordered input-sm px-3 py-1.5 rounded text-gray-700 w-full"
          />
        </div>

        <div class="mt-5">
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.greetings.message.title') }}</span>
          </div>
          <Field
            v-model.trim="greeting.text"
            name="text"
            as="input" 
            type="text"
            :placeholder="t('pages.greetings.message.placeholder')"
            :disabled="!greeting.edit"
            class="form-control input input-bordered input-sm px-3 py-1.5 rounded text-gray-700 w-full"
          />
        </div>
      </div>

      <div class="flex justify-between mt-5">
        <div>
          <MyBtn
            v-if="!greeting.edit"
            color="purple"
            @click="() => {
              greeting.edit = true;
              if (greeting.id) greetingsBeforeEdit?.push(JSON.parse(JSON.stringify(greeting)))
            }"
          >
            {{ t('buttons.edit') }}
          </MyBtn>
          <MyBtn
            v-else
            color="purple"
            @click="cancelEdit"
          >
            {{ t('buttons.cancel') }}
          </MyBtn>
        </div>
        <div
          v-if="greeting.edit"
          class="flex md:flex-none ml-1"
        >
          <MyBtn
            v-if="greeting.id"
            color="red"
            @click="deleteGreeting"
          >
            {{ t('buttons.delete') }}
          </MyBtn>
          <MyBtn
            color="green"
            type="submit"
            class="ml-1"
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