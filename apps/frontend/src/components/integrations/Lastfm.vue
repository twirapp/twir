<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ChannelIntegration } from '@tsuwari/prisma';
import {  ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

type Lastfm = Omit<ChannelIntegration, 'data'> & { data: { username: string }}

const lastfmIntegration = ref<Partial<Lastfm>>({
  enabled: true,
  data: {
    username: '',
  },
});
const selectedDashboard = useStore(selectedDashboardStore);
const { t } = useI18n({
  useScope: 'global',
});

selectedDashboardStore.subscribe(d => {
  api(`/v1/channels/${d.channelId}/integrations/lastfm`).then(async (r) => {
    if (r.data) {
      lastfmIntegration.value = r.data;
    } else {
      lastfmIntegration.value = {
        enabled: true,
        data: {
          username: '',
        },
      };
    }
  });
});

async function post() {
  const { data } = await api.post(`v1/channels/${selectedDashboard.value.channelId}/integrations/lastfm`, {
    enabled: lastfmIntegration.value.enabled,
    data: lastfmIntegration.value.data,
  });

  lastfmIntegration.value = data;
}
</script>

<template>
  <div class="card rounded card-compact bg-base-200 drop-shadow-lg p-4">
    <div class="flex justify-between mb-5">
      <div>
        <h2 class="card-title">
          Last.fm
        </h2>
      </div>
      <div class="form-check form-switch">
        <input
          id="flexSwitchCheckDefault"
          v-model="lastfmIntegration.enabled"
          class="form-check-input appearance-none w-9 -ml-10 rounded-full float-left h-5 align-top bg-no-repeat bg-contain bg-gray-300 focus:outline-none cursor-pointer shadow-sm"
          type="checkbox"
          role="switch"
        >
      </div>
    </div>

    <div>
      <div class="label">
        <span class="label-text">{{ t('pages.integrations.widgets.lastFm.username') }}</span>
      </div>
      <input
        v-model="lastfmIntegration.data!.username"
        type="text"
        class="form-control input text-gray-700 w-full flex-shrink flex-grow leading-normal rounded flex-1 border h-8 border-grey-light px-3 relative"
        :placeholder="t('pages.integrations.widgets.lastFm.username')"
      >
    </div>

    <button
      class="px-6 mt-3 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-purple-700 hover:shadow-lg focus:bg-purple-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-purple-800 active:shadow-lg transition duration-150 ease-in-out"
      @click="post"
    >
      {{ t('buttons.save') }}
    </button>
  </div>
</template>