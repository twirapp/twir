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
  <div class="m-1.5 md:m-3">
    <div class="flow-root">
      <div class="btn btn-primary btn-sm float-left mb-1 md:w-auto rounded w-full">
        <button
          class="bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-purple-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
          @click="insert"
        >
          {{ t('buttons.addNew') }}
        </button>
    

      <!-- <input
        type="text"
        placeholder="Search by keyword..."
        class="float-right rounded input input-sm input-bordered w-full md:w-60"
      > -->
      </div>
    </div>

    <div class="lg:masonry-lg masonry md:masonry-md sm:masonry-sm">
      <div
        v-for="timer, timerIndex of timers"
        :key="timerIndex"
        class="block break-inside card mb-[0.5rem] rounded shadow text-white"
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