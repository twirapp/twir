<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ChannelIntegration } from '@tsuwari/typeorm/entities/ChannelIntegration';
import { onMounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import { useToast } from 'vue-toastification';

import Tooltip from '../Tooltip.vue';

import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

type Integration = ChannelIntegration & {
  data: {
    name: string;
    code: string;
    avatar: string;
  };
};

const router = useRouter();
const donationAlertsIntegration = ref<Partial<Integration>>({
  enabled: true,
});
const selectedDashboard = useStore(selectedDashboardStore);

const { t } = useI18n({
  useScope: 'global',
});
const toast = useToast();

selectedDashboardStore.subscribe((d) => {
  api(`/v1/channels/${d.channelId}/integrations/donationalerts`).then(async (r) => {
    donationAlertsIntegration.value = r.data;
  });
});

async function auth() {
  const { data } = await api(
    `/v1/channels/${selectedDashboard.value.channelId}/integrations/donationalerts/auth`,
  );

  window.location.replace(data);
}

async function patch() {
  const { data } = await api.patch(
    `v1/channels/${selectedDashboard.value.channelId}/integrations/donationalerts`,
    {
      enabled: donationAlertsIntegration.value.enabled,
    },
  );

  donationAlertsIntegration.value = data;

  toast.success('Saved');
}

onMounted(async () => {
  const route = router.currentRoute.value;
  const params = new URLSearchParams(window.location.search);
  const code = params.get('code');

  if (route.params.integration === 'donationalerts' && code) {
    if (code) {
      await api.post(
        `v1/channels/${selectedDashboard.value.channelId}/integrations/donationalerts/token`,
        {
          code,
        },
      );

      return router.push('/dashboard/integrations');
    }
  }
});
</script>

<template>
  <div
    class="bg-base-200 break-inside card card-compact drop-shadow flex flex-col mb-[0.5rem] p-2 rounded">
    <div class="flex justify-between mb-5">
      <div>
        <h2 class="card-title flex font-bold space-x-2">
          <p>DonationAlerts</p>
          <Tooltip :text="t('pages.integrations.widgets.donationalerts.description')" />
        </h2>
      </div>
      <div class="form-check form-switch">
        <input
          v-model="donationAlertsIntegration.enabled"
          class="-ml-10 align-top appearance-none bg-contain bg-gray-300 bg-no-repeat cursor-pointer float-left focus:outline-none form-check-input h-5 rounded-full shadow w-9"
          type="checkbox"
          role="switch"
          @change="patch" />
      </div>
    </div>

    <div class="mb-5">
      <div v-if="donationAlertsIntegration.data">
        <div class="flex justify-center mb-3">
          <img
            v-if="donationAlertsIntegration.data.avatar"
            :src="donationAlertsIntegration.data.avatar"
            class="ring-2 ring-white rounded-full select-none w-32"
            alt="Avatar" />
        </div>
        <p class="break-words text-center">
          {{ donationAlertsIntegration.data.name }}#{{ donationAlertsIntegration.data.code }}
        </p>
      </div>
      <div v-else>Not logged in</div>
    </div>

    <div class="mt-auto text-right">
      <button
        class="bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-purple-700 inline-block leading-tight px-6 py-2.5 rounded shadow text-white text-xs transition uppercase"
        @click="auth">
        {{ t('buttons.login') }}
      </button>
    </div>
  </div>
</template>

<style scoped>
.tooltip {
  width: 140px;
  background: #59c7f9;
  color: #ffffff;
  text-align: center;
  padding: 10px 20px 10px 20px;
  border-radius: 10px;
  top: calc(100% + 11px);
  left: 50%;
  transform: translate-x(-50%);
}
.tooltip-box {
  position: relative;
}
.triangle {
  border-width: 0 6px 6px;
  border-color: transparent;
  border-bottom-color: #59c7f9;
  position: absolute;
  top: -6px;
  left: 50%;
  transform: translate-x(-50%);
}
</style>
