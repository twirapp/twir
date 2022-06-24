<script lang="ts" setup>import { useStore } from '@nanostores/vue';
import { useI18n } from 'vue-i18n';

import { api } from '@/plugins/api.js';
import { router } from '@/plugins/router.js';
import { selectedDashboardStore, userStore, setSelectedDashboard } from '@/stores/userStore';
import { setUser } from '@/stores/userStore.js';

const user = useStore(userStore);
const selectedDashboard = useStore(selectedDashboardStore);

async function logOut() {
  await api.post('/auth/logout');
  localStorage.clear();
  setUser(null);

  router.push('/');
}

const { t } = useI18n({
  useScope: 'global',
  inheritLocale: true,
});
</script>

<template>
  <!-- <img
    src="https://mdbcdn.b-cdn.net/img/new/avatars/2.webp"
    class="rounded-full w-9"
    alt="Avatar"
  > 

  <div class="dropdown relative">
    <img
      id="avatarDropdown"
      src="https://mdbcdn.b-cdn.net/img/new/avatars/2.webp"
      class="rounded-full w-9"
      alt="Avatar"
      data-bs-toggle="dropdown"
      aria-expanded="false"
    >

    <ul aria-labelledby="avatarDropdown">
      <p>sasi</p>
    </ul>
  </div> -->

  <div class="flex justify-center">
    <div>
      <div class="dropdown relative">
        <div class="absolute inline-block top-auto right-0 bottom-0 left-auto translate-x-1/4 translate-y-1/3 rotate-0 skew-x-0 skew-y-0 scale-x-100 scale-y-100 p-1.5 text-xs rounded-full z-10">
          <img
            v-if="selectedDashboard.channelId !== user?.id"
            id="profileMenu"
            type="button"
            data-bs-toggle="dropdown"
            aria-expanded="false"
            :src="user?.profile_image_url"
            class="
          rounded-full
          hover:cursor-pointer"
          
            alt="Avatar"
          > 
        </div>
        <img
          id="profileMenu"
          type="button"
          data-bs-toggle="dropdown"
          aria-expanded="false"
          :src="selectedDashboard?.twitch?.profile_image_url ?? user?.profile_image_url"
          class="
          rounded-full
          w-9
          hover:cursor-pointer"
          
          alt="Avatar"
        > 


        <div
          class="
          dropdown-menu
          min-w-max
          absolute
          w-64
          px-2
          bg-gray-700
          text-base
          z-50
          float-left
          py-2
          list-none
          text-left
          rounded
          mt-1
          hidden
          m-0
          bg-clip-padding
          border-none
        "
          aria-labelledby="profileMenu"
        >
          <div class="my-2 space-y-0.5">
            <span
              v-for="dashboard of user?.dashboards"
              :key="dashboard.channelId"
              :class="{'btn-disabled': selectedDashboard.channelId === dashboard.channelId}"
              class="dropdown-item
              text-sm
              py-2
              px-4
              font-normal
              block
              w-full
              whitespace-nowrap
              bg-transparent
              hover:bg-gray-600
              rounded
              cursor-pointer
            "
              @click="setSelectedDashboard(dashboard)"
            >
              <img
                class="w-6 rounded-full inline"
                :src="dashboard?.twitch?.profile_image_url ?? dashboard.twitch?.profile_image_url"
              >
              <span class="ml-4">{{ dashboard.twitch.display_name }}        
                <span
                  v-if="selectedDashboard?.id === user?.id"
                  class="text-xs inline-block py-1 px-2.5 leading-none text-center whitespace-nowrap align-baseline font-bold bg-gray-200 text-gray-700 rounded"
                >{{ t('pages.settings.widgets.dashboardAccess.thatsYou') }}</span>
              </span>
            </span>
          </div>
        
          <button
            class="inline-block w-full border-2 border-red-600 py-1.5 leading-tight uppercase rounded hover:bg-red-200 hover:bg-opacity-5 focus:outline-none focus:ring-0 transition duration-150 ease-in-out"
            @click="logOut"
          >
            Logout
          </button>
        </div>
      </div>
    </div>
  </div>
</template>