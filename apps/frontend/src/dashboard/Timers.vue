<script lang="ts" setup>
import { Ref, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import Timer from '@/components/Timer.vue';
import { useUpdatingData } from '@/functions/useUpdatingData';


const { t } = useI18n({
  useScope: 'global',
});
const timers = ref<Array<any>>([]);
const timersBeforeEdit = ref<Array<any>>([]);

const { data } = useUpdatingData(`/v1/channels/{dashboardId}/timers`);

watch(data, (v: any[]) => {
  timers.value = v;
  timersBeforeEdit.value = [];
});

function insert() {
  timers.value = [
    {
      name: '',
      enabled: true,
      last: 0,
      timeInterval: 60,
      messageInterval: 0,
      responses: [],
      edit: true,
    },
    ...timers.value,
  ];
}

function deleteTimer(index: number) {
  timers.value = timers.value.filter((_, i) => i !== index);
}

function cancelEdit(timer: Ref<any>) {
  const index = timers.value.indexOf(timer.value);
  if (timer.value.id && timers.value) {
    const editableCommand = timersBeforeEdit.value?.find(c => c.id === timer.value.id);
    if (editableCommand) {
      timers.value[index] = {
        ...editableCommand,
        edit: false,
      };

      timersBeforeEdit.value = timersBeforeEdit.value.filter((v, i) => i !== timersBeforeEdit.value.indexOf(editableCommand));
    }
  } else {
    timers.value = timers.value.filter((v, i) => i !== index);
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
      :items="timers"
      :gap="8"
    >
      <template #default="{ item, index }">
        <div
          :key="index"
          class="block card rounded shadow text-white"
        >
          <Timer
            :timer="item"
            :timers="timers"
            :timers-before-edit="timersBeforeEdit"
            @delete="deleteTimer"
            @cancel-edit="cancelEdit"
          />
        </div>
      </template>
    </masonry-wall>
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