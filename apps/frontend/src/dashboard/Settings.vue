<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { HelixUserData } from '@twurple/api';
import { ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import Add from '@/assets/buttons/add.svg';
import Remove from '@/assets/buttons/remove.svg';
import { useUpdatingData } from '@/functions/useUpdatingData';
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
const { data, execute } = useUpdatingData(`/v1/channels/{dashboardId}/settings/dashboardAccess`);

watch(data, (v) => {
  dashboardMembers.value = v;
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
  <div class="m-1.5 md:m-3">
    <div class="grid grid-cols-1 md:grid-cols-2">
      <div
        class="block card mb-[0.5rem] pb-0.5 rounded shadow text-white"
      >
        <h2 class="border-b border-gray-700 card-title flex font-bold justify-center outline-none p-2">
          <p>{{ t('pages.settings.widgets.dashboardAccess.title') }}</p>
        </h2>
        <div>
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
                    <img
                      class="inline rounded-full w-6"
                      :src="member.profile_image_url"
                    >
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
        </div>
      </div>
    </div>
  </div>
</template>
