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
    edit: true,
  });
}

async function deleteVariable(index: number) {
  variables.value = variables.value.filter((_, i) => i !== index);
}
</script>

<template>
  <div class="m-1.5 md:m-3">
    <div class="flow-root">
      <div class="btn btn-primary btn-sm float-left mb-1 md:w-auto rounded w-full">
        <button
          class="bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-purple-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
          @click="insert"
        >
          {{ t('buttons.addNew') }}
        </button>
      </div>
    </div>

    <div class="masonry">
      <div
        v-for="variable, index of variables"
        :key="index"
        class="block break-inside card mb-[0.5rem] rounded shadow text-white"
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
