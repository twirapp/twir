<script lang="ts" setup>
export type GreeTingType = SetOptional<Omit<Greeting, 'channelId'> & { username: string, edit?: boolean }, 'id'>

import { useStore } from '@nanostores/vue';
import { Greeting } from '@tsuwari/prisma';
import { useTitle } from '@vueuse/core';
import { useAxios } from '@vueuse/integrations/useAxios';
import type { SetOptional } from 'type-fest';
import { ref, watch } from 'vue';

import GreetingComponent from '@/components/Greeting.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const title = useTitle();
title.value = 'Tsuwari - Greetings';

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
  <div class="p-1">
    <div class="flow-root">
      <div class="float-left rounded btn btn-primary btn-sm w-full mb-1 md:w-auto">
        <button
          class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-purple-700 hover:shadow-lg focus:bg-purple-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out"
          @click="insert"
        >
          Add new
        </button>
      </div>
    </div>
  </div>

  <div class="grid lg:grid-cols-3 md:grid-cols-2 grid-cols-1 gap-2">
    <div
      v-for="greeting, index of greetings"
      :key="index"
      class="block rounded-lg card text-white shadow-lg"
    >
      <GreetingComponent 
        :greeting="greeting"
        :greetings="greetings"
        :greetings-before-edit="greetingsBeforeEdit"
        @delete="deleteGreeting"
      />
    </div>
  </div>
</template>
