<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ModerationSettingsDto } from '@tsuwari/shared';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import ModerationComponent from '@/components/Moderation.vue';
import { useUpdatingData } from '@/functions/useUpdatingData';
import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const settings = ref<ModerationSettingsDto[]>();
const { t } = useI18n({
  useScope: 'global',
});


const { data } = useUpdatingData(`/v1/channels/{dashboardId}/moderation`);

watch(data, (v) => {
  settings.value = v;
});

const selectedDashboard = useStore(selectedDashboardStore);

async function save() {
  await api.post(`/v1/channels/${selectedDashboard.value.channelId}/moderation`, settings.value);
}
</script>


<template>
  <div class="m-1.5 md:m-3">
    <div class="flow-root">
      <div class="btn btn-primary btn-sm float-left mb-1 md:w-auto rounded w-full">
        <button
          class="bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-purple-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
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
      class="lg:masonry-lg masonry md:masonry-md sm:masonry-sm"
    >
      <div
        v-for="setting, index in settings"
        :key="index"
        class="block break-inside card rounded shadow text-white"
      >
        <ModerationComponent :settings="(setting as any)" />
      </div>
    </div>
  </div>
</template>