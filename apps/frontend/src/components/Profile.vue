<script lang="ts" setup>
import MyBtn from '@elements/MyBtn.vue';
import { useStore } from '@nanostores/vue';
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';

import { api } from '@/plugins/api';
import { selectedDashboardStore, userStore, setSelectedDashboard } from '@/stores/userStore';
import { setUser } from '@/stores/userStore';

const user = useStore(userStore);
const selectedDashboard = useStore(selectedDashboardStore);
const router = useRouter();

const searchFilter = ref<string>('');
const filteredDashboards = computed(() => {
  return user.value?.dashboards
    .filter(c => searchFilter.value ? [c.twitch.login, c.twitch.display_name].some(s => s.includes(searchFilter.value)) : true)
    .sort((a, b) => a.twitch.id === user.value?.id ? -1 : a.twitch.login.localeCompare(b.twitch.login));
});

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
  <div class="flex justify-center select-none">
    <div>
      <div class="dropdown relative">
        <div
          id="profileMenu"
          type="button"
          data-bs-toggle="dropdown"
          aria-expanded="false"
          class="absolute bottom-0 hover:opacity-80 inline-block left-auto p-1.5 right-0 rotate-0 rounded-full scale-x-100 scale-y-100 skew-x-0 skew-y-0 text-xs top-auto translate-x-1/4 translate-y-1/3 z-10"
        >
          <img
            v-if="selectedDashboard.channelId !== user?.id"

            :src="user?.profile_image_url"
            class="hover:cursor-pointer
          rounded-full"
          
            alt="Avatar"
          > 
        </div>
        <img
          id="profileMenu"
          type="button"
          data-bs-toggle="dropdown"
          aria-expanded="false"
          :src="selectedDashboard?.twitch?.profile_image_url ?? user?.profile_image_url"
          class="hover:cursor-pointer
          hover:opacity-60
          rounded-full
          w-9"
          alt="Avatar"
        > 


        <div
          class="absolute
          bg-[#202020]
          bg-clip-padding
          border-none
          dropdown-menu
          float-left
          hidden
          list-none
          m-0
          min-w-max
          mt-1
          px-2
          py-2
          rounded
          text-base
          text-left
          w-64
          z-50"
          aria-labelledby="profileMenu"
        >
          <div
            v-if="!router.currentRoute.value.fullPath.startsWith('/admin')"
            class="max-h-[55vh] mb-2 overflow-auto overflow-y-auto scrollbar scrollbar-thin scrollbar-thumb-gray-600 scrollbar-track-gray-500 space-y-0.5"
          >
            <input
              v-if="(user?.dashboards && user.dashboards.length > 5) || user?.isBotAdmin"
              v-model="searchFilter"
              type="text"
              class="bg-clip-padding
              bg-white
              block
              border
              border-gray-300
              border-solid
              ease-in-out
              focus:bg-white
              focus:border-blue-600
              focus:outline-none
              focus:text-gray-700
              font-normal
              form-control
              m-0
              px-3
              py-0.5
              rounded
              text-base
              text-gray-700
              transition
              w-full"
              placeholder="Search..."
            >
            <div
              v-for="dashboard of filteredDashboards"
              :key="dashboard.id"
              :class="{'mr-2': user?.dashboards && user.dashboards.length > 5}"
            >
              <span
                :class="{
                  'border-2  border-cyan-600': dashboard.channelId === user?.id
                }"
                class="bg-transparent block cursor-pointer font-normal hover:bg-[#393636] mb-1 mt-1 px-4 py-2 rounded text-sm w-full whitespace-nowrap"
                @click="setSelectedDashboard(dashboard)"
              >
                <img
                  class="border inline rounded-full w-6"
                  :src="dashboard?.twitch?.profile_image_url ?? dashboard.twitch?.profile_image_url"
                >
                <span class="ml-4">{{ dashboard.twitch.display_name }}        
                  <span
                    v-if="selectedDashboard?.id === user?.id"
                    class="align-baseline bg-gray-200 font-bold inline-block leading-none px-2.5 py-1 rounded text-center text-gray-700 text-xs whitespace-nowrap"
                  >{{ t('pages.settings.widgets.dashboardAccess.thatsYou') }}</span>
                </span>
              </span>
            </div>
          </div>
        
          <div class="flex flex-col space-y-1 w-full">
            <MyBtn
              v-if="router.currentRoute.value.fullPath.startsWith('/admin')"
              color="purple"
              @click="router.push('/dashboard')"
            >
              Dashboard
            </MyBtn>
            <MyBtn
              v-if="!router.currentRoute.value.fullPath.startsWith('/admin') && user?.isBotAdmin"
              color="purple"
              @click="router.push('/admin')"
            >
              Admin
            </MyBtn>
            <MyBtn
              color="red"
              @click="logOut"
            >
              Logout
            </MyBtn>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.switch-button {
  @apply bg-purple-600 duration-150 ease-in-out focus:outline-none focus:ring-0 font-medium hover:bg-purple-700 inline-block leading-tight px-6 py-2 rounded text-sm text-white transition uppercase
}

</style>