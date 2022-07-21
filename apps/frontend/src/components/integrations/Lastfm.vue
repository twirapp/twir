<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ChannelIntegration } from '@tsuwari/prisma';
import {  ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useToast } from 'vue-toastification';

import Tooltip from '../Tooltip.vue';

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
const toast = useToast();

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
  toast.success('Saved');
}
</script>

<template>
  <div class="bg-base-200 break-inside card card-compact drop-shadow flex flex-col mb-[0.5rem] p-2 rounded">
    <div class="flex justify-between mb-5">
      <div>
        <h2 class="card-title flex font-bold space-x-2">
          <p>Last.fm</p>
          <Tooltip :text="t('pages.integrations.widgets.lastFm.description')" />
        </h2>
      </div>
      <div class="form-check form-switch">
        <input
          id="flexSwitchCheckDefault"
          v-model="lastfmIntegration.enabled"
          class="-ml-10 align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
          type="checkbox"
          role="switch"
        >
      </div>
    </div>

    <div>
      <div class="label mb-1">
        <span class="label-text">{{ t('pages.integrations.widgets.lastFm.username') }}</span>
      </div>
      <input
        v-model="lastfmIntegration.data!.username"
        type="text"
        class="border border-grey-light flex-1 flex-grow flex-shrink form-control h-8 input leading-normal px-3 relative rounded text-gray-700 w-full"
        :placeholder="t('pages.integrations.widgets.lastFm.username')"
      >
    </div>

 
    <div class="mt-auto text-right">
      <button
        class="bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-purple-700 inline-block leading-tight mt-3 px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
        @click="post"
      >
        {{ t('buttons.save') }}
      </button>
    </div>
  </div>
</template>