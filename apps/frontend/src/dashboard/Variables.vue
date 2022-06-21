<script lang="ts" setup>
export type VariableType = SetOptional<Omit<CustomVar, 'channelId'| 'evalValue'> & { edit?: boolean, evalValue: string }, 'id'>

import { useStore } from '@nanostores/vue';
import { CustomVar } from '@tsuwari/prisma';
import { useAxios } from '@vueuse/integrations/useAxios';
import type { SetOptional } from 'type-fest';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import VariableComponent from '@/components/Variable.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const selectedDashboard = useStore(selectedDashboardStore);

const { execute, data: axiosData } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/variables`, api, { immediate: false });
const variables = ref<Array<VariableType>>([]);
const variablesBeforeEdit = ref<Array<VariableType>>([]);
const { t } = useI18n({
  useScope: 'global',
});

selectedDashboardStore.subscribe((v) => {
  execute(`/v1/channels/${v.channelId}/variables`);
});

watch(axiosData, (v: any[]) => {
  variables.value = v;
  variablesBeforeEdit.value = [];
});

function insert() {
  variables.value.unshift({
    name: '',
    description: '',
    type: 'SCRIPT',
    evalValue: '',
    response: null,
  });
}

async function deleteVariable(index: number) {
  variables.value = variables.value.filter((_, i) => i !== index);
}
</script>

<template>
  <div class="m-3">
    <div class="p-1">
      <div class="flow-root">
        <div class="float-left rounded btn btn-primary btn-sm w-full mb-1 md:w-auto">
          <button
            class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-purple-700 hover:shadow-lg focus:bg-purple-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out"
            @click="insert"
          >
            {{ t('buttons.addNew') }}
          </button>
        </div>
      </div>
    </div>

    <div class="grid lg:grid-cols-2 md:grid-cols-1 grid-cols-1 gap-2">
      <div
        v-for="variable, index of variables"
        :key="index"
        class="block rounded-lg card text-white shadow-lg"
      >
        <VariableComponent 
          :variable="variable"
          :variables="variables"
          :variables-before-edit="variablesBeforeEdit"
          @delete="deleteVariable"
        />
      </div>
    </div>
  </div>
</template>
