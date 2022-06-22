<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { useAxios } from '@vueuse/integrations/useAxios';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import Timer from '@/components/Timer.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';


const { t } = useI18n({
  useScope: 'global',
});
const selectedDashboard = useStore(selectedDashboardStore);
const timers = ref<Array<any>>([]);
const timersBeforeEdit = ref<Array<any>>([]);

const { execute, data: axiosData } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/timers`, api, { immediate: false });

selectedDashboardStore.subscribe((v) => {
  execute(`/v1/channels/${v.channelId}/timers`);
});

watch(axiosData, (v: any[]) => {
  timers.value = v;
  timersBeforeEdit.value = [];
});

function insert() {
  timers.value.unshift({
    name: '',
    enabled: true,
    last: 0,
    timeInterval: 60,
    messageInterval: 0,
    responses: [],
    edit: true,
  });
}

function deleteTimer(index: number) {
  timers.value = timers.value.filter((_, i) => i !== index);
}
</script>

<template>
  <div class="m-3">
    <div class="p-1">
      <div class="flow-root">
        <div class="float-left rounded btn btn-primary btn-sm w-full mb-1 md:w-auto">
          <button
            class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-purple-700 hover:shadow focus:bg-purple-700 focus:shadow focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow transition duration-150 ease-in-out"
            @click="insert"
          >
            {{ t('buttons.addNew') }}
          </button>
        </div>

      <!-- <input
        type="text"
        placeholder="Search by keyword..."
        class="float-right rounded input input-sm input-bordered w-full md:w-60"
      > -->
      </div>
    </div>

    <div class="grid xl:grid-cols-3 lg:grid-cols-2 grid-cols-1 gap-2">
      <div
        v-for="timer, timerIndex of timers"
        :key="timerIndex"
        class="block rounded card text-white shadow"
      >
        <Timer
          :timer="timer"
          :timers="timers"
          :timers-before-edit="timersBeforeEdit"
          @delete="deleteTimer"
        />
      </div>
    </div>
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