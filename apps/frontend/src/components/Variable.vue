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

import MyBtn from './elements/MyBtn.vue';

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
        class="bg-red-600 mb-3 mb-4 px-6 py-5 rounded text-base text-red-700"
        role="alert"
      >
        {{ error }}
      </div>
      <div
        class="gap-1 grid grid-cols-1"
      >
        <div>
          <div class="label mb-1">
            <span class="label-text">{{ t('pages.variables.card.name.title') }}</span>
          </div>
          <Field
            v-model="variable.name"
            name="name"
            as="input" 
            type="text"
            :placeholder="t('pages.variables.card.name.placeholder')"
            :disabled="!variable.edit"
            class="form-control input input-bordered input-sm px-3 py-1.5 rounded text-gray-700 w-full"
          />
        </div>
      </div>

      <div class="mt-5">
        <div class="label mb-1 mt-2">
          <span class="label-text">{{ t('pages.variables.card.type') }}</span>
        </div>
        <select
          v-model="variable.type"
          :disabled="!variable.edit"
          class="form-control px-3 py-1.5 rounded select select-sm text-gray-700 w-full"
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
        class="mt-5"
      >
        <prism-editor
          v-model="variable.evalValue"
          class="my-editor"
          :highlight="highlighter"
          :line-numbers="loaded"
          :readonly="!variable.edit"
        />
      </div>

      <div
        v-if="variable.type === 'TEXT'"
        class="mt-5"
      >
        <div class="label mb-1">
          <span class="label-text">{{ t('pages.variables.card.messageForSending.title') }}</span>
        </div>
        <Field
          v-model="variable.response"
          name="response"
          as="input" 
          type="text"
          :placeholder="t('pages.variables.card.messageForSending.placeholder')"
          :disabled="!variable.edit"
          class="form-control input input-bordered input-sm px-3 py-1.5 rounded text-gray-700 w-full"
        />
      </div>

      <div class="flex justify-between mt-5">
        <div>
          <MyBtn
            v-if="!variable.edit"
            color="purple"
            @click="() => {
              variable.edit = true;
              if (variable.id) variablesBeforeEdit?.push(JSON.parse(JSON.stringify(variable)))
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
        <div v-if="variable.edit">
          <MyBtn
            v-if="variable.id"
            color="red"
            @click="deleteVariable"
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