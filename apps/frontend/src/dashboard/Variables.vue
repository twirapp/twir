<script lang="ts" setup>
export type VariableType = SetOptional<Omit<CustomVar, 'channelId'| 'evalValue'> & { edit?: boolean, evalValue: string }, 'id'>

import { CustomVar } from '@tsuwari/prisma';
import type { SetOptional } from 'type-fest';
import { ref, watch, Ref } from 'vue';
import { useI18n } from 'vue-i18n';

import VariableComponent from '@/components/Variable.vue';
import { useUpdatingData } from '@/functions/useUpdatingData';


const { data } = useUpdatingData(`/v1/channels/{dashboardId}/variables`);

const variables = ref<Array<VariableType>>([]);
const variablesBeforeEdit = ref<Array<VariableType>>([]);
const { t } = useI18n({
  useScope: 'global',
});

watch(data, (v: any[]) => {
  variables.value = v;
  variablesBeforeEdit.value = [];
});

function insert() {
  variables.value = [{
    name: '',
    description: '',
    type: 'SCRIPT',
    evalValue: '',
    response: null,
    edit: true,
  }, ...variables.value];
}

async function deleteVariable(index: number) {
  variables.value = variables.value.filter((_, i) => i !== index);
}

function cancelEdit(variable: Ref<VariableType>) {
  const index = variables.value.indexOf(variable.value);
  if (variable.value.id) {
    const editableCommand = variablesBeforeEdit.value?.find(c => c.id === variable.value.id);
    if (editableCommand) {
      variables.value[index] = {
        ...editableCommand,
        edit: false,
      };

      variablesBeforeEdit.value = variablesBeforeEdit.value.filter((v, i) => i !== variablesBeforeEdit.value.indexOf(editableCommand));
    }
  } else {
    variables.value = variables.value.filter((v, i) => i !== index);
  }
}
</script>

<template>
  <div class="m-1.5 md:m-3">
    <div class="flow-root">
      <div class="btn btn-primary btn-sm float-left mb-5 md:w-auto rounded w-full">
        <button
          class="bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-purple-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
          @click="insert"
        >
          {{ t('buttons.addNew') }}
        </button>
      </div>
    </div>

    <masonry-wall
      :items="variables"
      :gap="8"
    >
      <template #default="{ item, index }">
        <div
          :key="index"
          class="block card rounded shadow text-white"
        >
          <VariableComponent 
            :variable="item"
            :variables="variables"
            :variables-before-edit="variablesBeforeEdit"
            @delete="deleteVariable"
            @cancel-edit="cancelEdit"
          />
        </div>
      </template>
    </masonry-wall>
  </div>
</template>
