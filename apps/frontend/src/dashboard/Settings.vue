<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { HelixUserData } from '@twurple/api';
import { useAxios } from '@vueuse/integrations/useAxios';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import Add from '@/assets/buttons/add.svg';
import Remove from '@/assets/buttons/remove.svg';
import { api } from '@/plugins/api';
import { selectedDashboardStore, userStore } from '@/stores/userStore';

const selectedDashboard = useStore(selectedDashboardStore);
const { t } = useI18n({
  useScope: 'global',
  inheritLocale: true,
});
const user = useStore(userStore);
const newMember = ref('');
const dashboardMembers = ref<Array<HelixUserData>>();
const { execute, data: axiosData } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/settings/dashboardAccess`, api, { immediate: false });

watch(axiosData, (v) => {
  dashboardMembers.value = v;
});

selectedDashboardStore.subscribe(async (v) => {
  execute(`/v1/channels/${v.channelId}/settings/dashboardAccess`);
});

async function deleteMember(id: string) {
  await api.delete(`/v1/channels/${selectedDashboard.value.channelId}/settings/dashboardAccess/${id}`);
  execute();
}

async function addMember() {  
  await api.post(`/v1/channels/${selectedDashboard.value.channelId}/settings/dashboardAccess`, { username: newMember.value });
  execute();
  newMember.value = '';
}
</script>

<template>
  <div class="p-1 m-3">
    <div class="grid lg:grid-cols-3 grid-cols-1 gap-2">
      <div
        class="block rounded card text-white shadow max-w-sm"
      >
        <h2 class="card-title p-2 flex justify-center border-b border-gray-700 outline-none font-bold">
          <p>{{ t('pages.settings.widgets.dashboardAccess.title') }}</p>
        </h2>
        <div>
          <div
            class="rounded text-base my-4"
          >
            <p
              v-if="!dashboardMembers?.length"
              class="text-center"
            >
              {{ t('pages.settings.widgets.dashboardAccess.noAccs') }}
            </p>
            <ul
              v-else
              class="w-full max-h-[55vh] overflow-y-auto scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600"
            >
              <li
                v-for="member of dashboardMembers"
                :key="member.id"
                class="
                  text-sm
                  py-2
                  px-4
                  font-normal
                  block
                  w-full
                  whitespace-nowrap
                  bg-transparent"
                :class="{
                  'hover:bg-[#121212]': member.id !== user?.id,
                  'bg-[#121212]': member.id === user?.id
                }"
              >
                <div
                  class="flex justify-between items-center"
                  :class="{'cursor-pointer': member.id !== user?.id}"
                  @click="member.id !== user?.id ? deleteMember(member.id) : null"
                >
                  <div>
                    <img
                      class="w-6 rounded-full inline"
                      :src="member.profile_image_url"
                    >
                    <span class="ml-4">{{ member.display_name }}</span>
                  </div>
                  <div v-if="member.id !== user?.id">
                    <Remove />
                  </div>
                  <div v-if="member.id === user?.id">
                    <span class="text-xs inline-block py-1 px-2.5 leading-none text-center whitespace-nowrap align-baseline font-bold bg-gray-200 text-gray-700 rounded">{{ t('pages.settings.widgets.dashboardAccess.thatsYou') }}</span>
                  </div>
                </div>
              </li>
              <ul />
            </ul>
          </div>
          <div class="flex flex-wrap items-stretch w-full relative">
            <input
              v-model="newMember"
              type="text"
              class="form-control rounded-l input text-gray-700 flex-shrink flex-grow leading-normal flex-1 border h-10 border-grey-light px-3 relative"
              :placeholder="t('pages.settings.widgets.dashboardAccess.placeholder')"
              @keyup.enter="addMember"
            >
            <div
              class="flex cursor-pointer"
              @click="addMember"
            >
              <span class="flex items-center leading-normal bg-green-600 hover:bg-green-700 px-3 whitespace-no-wrap text-grey-dark text-sm rounded-r"><Add /></span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
