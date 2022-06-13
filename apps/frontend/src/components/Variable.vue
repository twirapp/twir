<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import type { CustomVarType as TCustomVarType } from '@tsuwari/prisma';
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
//@ts-ignore
import { highlight, languages } from 'prismjs/components/prism-core';
import { Form, Field } from 'vee-validate';
import { onMounted, ref, toRef } from 'vue';
import { useI18n } from 'vue-i18n';
import { PrismEditor } from 'vue-prism-editor';

import 'vue-prism-editor/dist/prismeditor.min.css';
import 'prismjs/components/prism-clike';
import 'prismjs/components/prism-javascript';
import 'prismjs/themes/prism-tomorrow.css';

import { VariableType } from '@/dashboard/Variables.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const CustomVarType = {
  'Script': 'SCRIPT',
  'Text': 'TEXT',
} as { [x: string]: TCustomVarType };
const highlighter = (code: string) => highlight(code, languages.js);

const props = defineProps<{
  variable: VariableType,
  variables: VariableType[],
  variablesBeforeEdit: VariableType[]
}>();

const variable = toRef(props, 'variable');
const variables = toRef(props, 'variables');
const variablesBeforeEdit = toRef(props, 'variablesBeforeEdit');
const selectedDashboard = useStore(selectedDashboardStore);
const { t } = useI18n({
  useScope: 'global',
});

const emit = defineEmits<{
  (e: 'delete', index: number): void
}>();

async function saveVariable() {
  const index = variables.value.indexOf(variable.value);
  
  let data;
  if (props.variable.id) {
    const request = await api.put(`/v1/channels/${selectedDashboard.value.channelId}/variables/${variable.value.id}`, variable.value);  
    data = request.data;
  } else {
    const request = await api.post(`/v1/channels/${selectedDashboard.value.channelId}/variables`, variable.value);
    data = request.data;
  }

  variables.value[index] = data;
}

async function deleteVariable() {
  const index = variables.value.indexOf(variable.value);
  if (variable.value.id) {
    await api.delete(`/v1/channels/${selectedDashboard.value.channelId}/variables/${variable.value.id}`);
  }

  emit('delete', index);
}

function cancelEdit() {
  const index = variables.value.indexOf(variable.value);
  if (variable.value.id && variables.value) {
    const editableCommand = variablesBeforeEdit.value?.find(c => c.id === variable.value.id);
    if (editableCommand) {
      variables.value[index] = {
        ...editableCommand,
        edit: false,
      };
      variablesBeforeEdit.value?.splice(variablesBeforeEdit.value.indexOf(editableCommand), 1);
    }
  } else {
    variables.value?.splice(index, 1);
  }
}
const loaded = ref(false);

onMounted(() => {
  setTimeout(() => (loaded.value = true), 100);
});

</script>

<template>
  <div class="p-4">
    <Form
      v-slot="{ errors }"
      @submit="saveVariable"
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
        class="grid grid-cols-1 gap-1"
      >
        <div>
          <div class="label mb-3">
            <span class="label-text">{{ t('pages.variables.card.name.title') }}</span>
          </div>
          <Field
            v-model="variable.name"
            name="name"
            as="input" 
            type="text"
            :placeholder="t('pages.variables.card.name.placeholder')"
            :disabled="!variable.edit"
            class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
          />
        </div>
      </div>

      <div>
        <div class="label mt-2 mb-3">
          <span class="label-text">{{ t('pages.variables.card.type') }}</span>
        </div>
        <select
          v-model="variable.type"
          :disabled="!variable.edit"
          class="form-control px-3 py-1.5 text-gray-700 rounded select select-sm w-full"
        >
          <option
            v-for="type of Object.entries(CustomVarType)"
            :key="type[0]"
            :value="type[1]"
          >
            {{ type[0] }}
          </option>
        </select>
      </div>

      <div
        v-if="variable.type === 'SCRIPT'"
        class="mt-3"
      >
        <prism-editor
          v-model="variable.evalValue"
          class="my-editor"
          :highlight="highlighter"
          :line-numbers="loaded"
          :readonly="!variable.edit"
        />
      </div>

      <div v-if="variable.type === 'TEXT'">
        <div class="label mb-3">
          <span class="label-text">{{ t('pages.variables.card.messageForSending.title') }}</span>
        </div>
        <Field
          v-model="variable.response"
          name="response"
          as="input" 
          type="text"
          :placeholder="t('pages.variables.card.messageForSending.placeholder')"
          :disabled="!variable.edit"
          class="form-control px-3 py-1.5 text-gray-700 rounded input input-bordered w-full input-sm"
        />
      </div>

      <div class="flex justify-between mt-5">
        <div>
          <button
            v-if="!variable.edit"
            type="button"
            class="inline-block px-6 py-2.5 bg-gray-200 text-gray-700 font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-gray-300 hover:shadow-lg focus:bg-gray-300 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-gray-400 active:shadow-lg transition duration-150 ease-in-out"
            @click="() => {
              variable.edit = true;
              if (variable.id) variablesBeforeEdit?.push(JSON.parse(JSON.stringify(variable)))
            }"
          >
            {{ t('buttons.edit') }}
          </button>
          <button
            v-else
            class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-purple-700 hover:shadow-lg focus:bg-purple-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out"
            @click="cancelEdit"
          >
            {{ t('buttons.cancel') }}
          </button>
        </div>
        <div v-if="variable.edit">
          <button
            v-if="variable.id"
            type="button"
            class="inline-block px-6 py-2.5 bg-red-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-red-700 hover:shadow-lg focus:bg-red-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-red-800 active:shadow-lg transition duration-150 ease-in-out"
            @click="deleteVariable"
          >
            {{ t('buttons.delete') }}
          </button>
          <button
            type="submit"
            class="inline-block ml-2 px-6 py-2.5 bg-green-500 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-green-600 hover:shadow-lg focus:bg-green-600 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-green-700 active:shadow-lg transition duration-150 ease-in-out"
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