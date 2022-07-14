<script lang="ts" setup>
export type GreeTingType = SetOptional<Omit<Greeting, 'channelId'> & { username: string, edit?: boolean }, 'id'>

import { useStore } from '@nanostores/vue';
import { Greeting } from '@tsuwari/prisma';
import { useAxios } from '@vueuse/integrations/useAxios';
import type { SetOptional } from 'type-fest';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import GreetingComponent from '@/components/Greeting.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const { t } = useI18n({
  useScope: 'global',
});

const selectedDashboard = useStore(selectedDashboardStore);

const { execute, data: axiosData } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/greetings`, api, { immediate: false });
const greetings = ref<Array<GreeTingType>>([]);
const greetingsBeforeEdit = ref<Array<GreeTingType>>([]);

selectedDashboardStore.subscribe((v) => {
  execute(`/v1/channels/${v.channelId}/greetings`);
});

watch(axiosData, (v: any[]) => {
  greetings.value = v;
  greetingsBeforeEdit.value = [];
});

function insert() {
  greetings.value.unshift({
    username: '',
    userId: '',
    text: '',
    edit: true,
    enabled: true,
  });
}

async function deleteGreeting(index: number) {
  greetings.value = greetings.value.filter((_, i) => i !== index);
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
          {{ t('pages.greetings.buttons.add') }}
        </button>
      </div>
    </div>


    <div class="gap-2 grid grid-cols-1 md:grid-cols-2">
      <div
        v-for="greeting, index of greetings"
        :key="index"
        class="block card rounded shadow text-white"
      >
        <GreetingComponent 
          :greeting="greeting"
          :greetings="greetings"
          :greetings-before-edit="greetingsBeforeEdit"
          @delete="deleteGreeting"
        />
      </div>
    </div>
  </div>
</template>
