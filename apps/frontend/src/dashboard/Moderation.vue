<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ModerationSettingsDto } from '@tsuwari/shared';
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';

import ModerationComponent from '@/components/Moderation.vue';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const settings = ref<ModerationSettingsDto[]>();
const { t } = useI18n({
  useScope: 'global',
});

const selectedDashboard = useStore(selectedDashboardStore);

selectedDashboardStore.subscribe(() => {
  getModerationSettings();
});

async function getModerationSettings() {
  const { data } = await api(`/v1/channels/${selectedDashboard.value.channelId}/moderation`);
  settings.value = data;
}

onMounted(() => {
  getModerationSettings();
});

async function save() {
  await api.post(`/v1/channels/${selectedDashboard.value.channelId}/moderation`, settings.value);
}
</script>


<template>
  <div class="m-1.5 md:m-3">
    <div class="flow-root">
      <div class="float-left rounded btn btn-primary btn-sm w-full mb-1 md:w-auto">
        <button
          class="px-6 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-purple-700    focus:outline-none focus:ring-0  transition duration-150 ease-in-out"
          @click="save"
        >
          {{ t('buttons.save') }}
        </button>


      <!-- <input
        type="text"
        placeholder="Search by keyword..."
        class="float-right rounded input input-sm input-bordered w-full md:w-60"
      > -->
      </div>
    </div>
    <div 
      class="masonry sm:masonry-sm md:masonry-md lg:masonry-lg"
    >
      <div
        v-for="setting, index in settings"
        :key="index"
        class="break-inside block rounded card text-white shadow"
      >
        <ModerationComponent :settings="(setting as any)" />
      </div>
    </div>
  </div>
</template>