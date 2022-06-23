<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ChannelIntegration } from '@tsuwari/prisma';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

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
}
</script>

<template>
  <div class="flex flex-col card rounded card-compact bg-base-200 drop-shadow p-2 break-inside mb-[0.5rem]">
    <div class="flex justify-between mb-5">
      <div>
        <h2 class="card-title font-bold">
          Faceit
        </h2>
      </div>
      <div class="form-check form-switch">
        <input
          id="flexSwitchCheckDefault"
          v-model="faceitIntegration.enabled"
          class="form-check-input appearance-none w-9 -ml-10 rounded-full float-left h-5 align-top bg-no-repeat bg-contain bg-gray-300 focus:outline-none cursor-pointer shadow"
          type="checkbox"
          role="switch"
        >
      </div>
    </div>

    <div>
      <div class="label">
        <span class="label-text">{{ t('pages.integrations.widgets.faceit.username') }}</span>
      </div>
      <input
        v-model="faceitIntegration.data!.username"
        type="text"
        class="form-control input text-gray-700 w-full flex-shrink flex-grow leading-normal rounded flex-1 border h-8 border-grey-light px-3 relative"
        :placeholder="t('pages.integrations.widgets.faceit.username')"
      >
    </div>

    <div>
      <div class="label">
        <span class="label-text">{{ t('pages.integrations.widgets.vk.id') }}</span>
      </div>
      <select
        v-model="faceitIntegration.data!.game"
        class="form-control px-3 py-1.5 text-gray-700 rounded select select-sm w-full"
      >
        <option value="csgo">
          CSGO
        </option>
      </select>
    </div>

    <div class="mt-auto text-right">
      <button
        class="px-6 mt-3 py-2.5 inline-block bg-purple-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-purple-700    focus:outline-none focus:ring-0  transition duration-150 ease-in-out"
        @click="post"
      >
        {{ t('buttons.save') }}
      </button>
    </div>
  </div>
</template>