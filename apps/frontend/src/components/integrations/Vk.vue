<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ChannelIntegration } from '@tsuwari/prisma';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';


type Vk = Omit<ChannelIntegration, 'data'> & { data: { userId: string }}

const vkIntegration = ref<Partial<Vk>>({
  enabled: true,
  data: {
    userId: '',
  },
});
const selectedDashboard = useStore(selectedDashboardStore);
const { t } = useI18n({
  useScope: 'global',
});

selectedDashboardStore.subscribe(d => {
  api(`/v1/channels/${d.channelId}/integrations/vk`).then(async (r) => {
    if (r.data) {
      vkIntegration.value = r.data;
    } else {
      vkIntegration.value = {
        enabled: true,
        data: {
          userId: '',
        },
      };
    }
  });
});

async function post() {
  const { data } = await api.post(`v1/channels/${selectedDashboard.value.channelId}/integrations/vk`, {
    enabled: vkIntegration.value.enabled,
    data: vkIntegration.value.data,
  });

  vkIntegration.value = data;
}
</script>

<template>
  <div class="flex flex-col card rounded card-compact bg-base-200 drop-shadow p-4">
    <div class="flex justify-between mb-5">
      <div>
        <h2 class="card-title font-bold">
          VK
        </h2>
      </div>
      <div class="form-check form-switch">
        <input
          id="flexSwitchCheckDefault"
          v-model="vkIntegration.enabled"
          class="form-check-input appearance-none w-9 -ml-10 rounded-full float-left h-5 align-top bg-no-repeat bg-contain bg-gray-300 focus:outline-none cursor-pointer shadow"
          type="checkbox"
          role="switch"
        >
      </div>
    </div>

    <div>
      <div class="label">
        <span class="label-text">{{ t('pages.integrations.widgets.vk.id') }}</span>
      </div>
      <input
        v-model="vkIntegration.data!.userId"
        type="text"
        class="form-control input text-gray-700 w-full flex-shrink flex-grow leading-normal rounded flex-1 border h-8 border-grey-light px-3 relative"
        :placeholder="t('pages.integrations.widgets.vk.id')"
      >
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