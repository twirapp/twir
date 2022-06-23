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
    <div class="p-1">
      <div class="flow-root">
        <div class="float-left rounded btn btn-primary btn-sm w-full mb-1 md:w-auto">
          <button
            class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-purple-700    focus:outline-none focus:ring-0  transition duration-150 ease-in-out"
            @click="insert"
          >
            {{ t('pages.greetings.buttons.add') }}
          </button>
        </div>
      </div>
    </div>

    <div class="masonry sm:masonry-sm md:masonry-md lg:masonry-lg">
      <div
        v-for="greeting, index of greetings"
        :key="index"
        class="block rounded card text-white shadow break-inside mb-[0.5rem]"
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
