<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ChannelIntegration } from '@tsuwari/prisma';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useToast } from 'vue-toastification';

import Tooltip from '../Tooltip.vue';

import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

type Faceit = Omit<ChannelIntegration, 'data'> & { data: { username: string, game?: string }}

const faceitIntegration = ref<Partial<Faceit>>({
  enabled: true,
  data: {
    username: '',
    game: 'csgo',
  },
});
const selectedDashboard = useStore(selectedDashboardStore);
const { t } = useI18n({
  useScope: 'global',
});
const toast = useToast();

selectedDashboardStore.subscribe(d => {
  api(`/v1/channels/${d.channelId}/integrations/faceit`).then(async (r) => {
    if (r.data) {
      faceitIntegration.value = r.data;
    } else {
      faceitIntegration.value = {
        enabled: true,
        data: {
          username: '',
          game: 'csgo',
        },
      };
    }
  });
});

async function post() {
  const { data } = await api.post(`v1/channels/${selectedDashboard.value.channelId}/integrations/faceit`, {
    enabled: faceitIntegration.value.enabled,
    data: faceitIntegration.value.data,
  });

  faceitIntegration.value = data;
  toast.success('Saved');
}
</script>

<template>
  <div class="bg-base-200 break-inside card card-compact drop-shadow flex flex-col mb-[0.5rem] p-2 rounded">
    <div class="flex justify-between mb-5">
      <div>
        <h2 class="card-title flex font-bold space-x-2">
          <p>FaceIT</p>
          <Tooltip
            :text="t('pages.integrations.widgets.faceit.description')"
          />
        </h2>
      </div>
      <div class="form-check form-switch">
        <input
          id="flexSwitchCheckDefault"
          v-model="faceitIntegration.enabled"
          class="-ml-10 align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
          type="checkbox"
          role="switch"
        >
      </div>
    </div>

    <div>
      <div class="label mb-1">
        <span class="label-text">{{ t('pages.integrations.widgets.faceit.username') }}</span>
      </div>
      <input
        v-model="faceitIntegration.data!.username"
        type="text"
        class="border border-grey-light flex-1 flex-grow flex-shrink form-control h-8 input leading-normal px-3 relative rounded text-gray-700 w-full"
        :placeholder="t('pages.integrations.widgets.faceit.username')"
      >
    </div>

    <div class="mt-5">
      <div class="label mb-1">
        <span class="label-text">{{ t('pages.integrations.widgets.faceit.game') }}</span>
      </div>
      <select
        v-model="faceitIntegration.data!.game"
        class="form-control px-3 py-1.5 rounded select select-sm text-gray-700 w-full"
      >
        <option value="csgo">
          CSGO
        </option>
      </select>
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