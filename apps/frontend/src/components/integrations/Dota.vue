<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { ChannelIntegration } from '@tsuwari/prisma';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useToast } from 'vue-toastification';

import Soon from '../Soon.vue';
import Tooltip from '../Tooltip.vue';

import Add from '@/assets/buttons/add.svg';
import Remove from '@/assets/buttons/remove.svg';
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
  const { data } = await api.post(`v1/channels/${selectedDashboard.value.channelId}/integrations/dota2`, {
    enabled: faceitIntegration.value.enabled,
    data: faceitIntegration.value.data,
  });

  faceitIntegration.value = data;
  toast.success('Saved');
}
</script>

<template>
  <div class="bg-base-200 break-inside card card-compact drop-shadow flex flex-col mb-[0.5rem] rounded">
    <div class="flex justify-between mb-5 p-2">
      <div>
        <h2 class="card-title flex font-bold space-x-2">
          <p>Dota 2</p>
          <Tooltip
            :text="t('pages.integrations.widgets.dota2.description')"
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


    <Soon :button="false" />
    
    <!-- <div>
      <div
        class="my-4 rounded text-base"
      >
        <p
          v-if="!dashboardMembers?.length"
          class="text-center"
        >
          {{ t('pages.settings.widgets.dashboardAccess.noAccs') }}
        </p>
        <ul
          v-else
          class="max-h-[55vh] overflow-auto overflow-y-auto scrollbar scrollbar-thin scrollbar-thumb-gray-900 scrollbar-track-gray-600 w-full"
        >
          <li
            v-for="member of dashboardMembers"
            :key="member.id"
            class="bg-transparent block font-normal px-4 py-2 text-sm w-full whitespace-nowrap"
            :class="{
              'hover:bg-[#121212]': member.id !== user?.id,
              'bg-[#121212]': member.id === user?.id
            }"
          >
            <div
              class="flex items-center justify-between"
              :class="{'cursor-pointer': member.id !== user?.id}"
              @click="member.id !== user?.id ? deleteMember(member.id) : null"
            >
              <div>
                <span class="ml-4">{{ member.display_name }}</span>
              </div>
              <div v-if="member.id !== user?.id">
                <Remove />
              </div>
              <div v-if="member.id === user?.id">
                <span class="align-baseline bg-gray-200 font-bold inline-block leading-none px-2.5 py-1 rounded text-center text-gray-700 text-xs whitespace-nowrap">{{ t('pages.settings.widgets.dashboardAccess.thatsYou') }}</span>
              </div>
            </div>
          </li>
          <ul />
        </ul>
      </div>
      <div class="flex flex-wrap items-stretch relative w-full">
        <input
          v-model="newMember"
          type="text"
          class="border border-grey-light flex-1 form-control h-10 input leading-normal px-3 relative rounded-l text-gray-700 w-full"
          :placeholder="t('pages.settings.widgets.dashboardAccess.placeholder')"
          @keyup.enter="addMember"
        >
        <div
          class="cursor-pointer flex"
          @click="addMember"
        >
          <span class="bg-green-600 flex hover:bg-green-700 items-center leading-normal px-3 rounded-r text-grey-dark text-sm whitespace-no-wrap"><Add /></span>
        </div>
      </div>
    </div> -->
  </div>
</template>