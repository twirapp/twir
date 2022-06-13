<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { useTitle } from '@vueuse/core';
import { useI18n } from 'vue-i18n';
import { useRouter, useRoute } from 'vue-router';

import Commands from '@/assets/sidebar/commands.svg?url';
import Dashboard from '@/assets/sidebar/dashboard.svg?url';
import Events from '@/assets/sidebar/events.svg?url';
import Files from '@/assets/sidebar/files.svg?url';
import Greetings from '@/assets/sidebar/greetings.svg?url';
import Integrations from '@/assets/sidebar/integrations.svg?url';
import Keywords from '@/assets/sidebar/keywords.svg?url';
import Overlays from '@/assets/sidebar/overlays.svg?url';
import Quotes from '@/assets/sidebar/quotes.svg?url';
import Settings from '@/assets/sidebar/settings.svg?url';
import Sword from '@/assets/sidebar/sword.svg?url';
import Timers from '@/assets/sidebar/timers.svg?url';
import Users from '@/assets/sidebar/users.svg?url';
import Variables from '@/assets/sidebar/variables.svg?url';
import { selectedDashboardStore, userStore, setSelectedDashboard } from '@/stores/userStore';

const user = useStore(userStore);
const router = useRouter();
const { t } = useI18n({
  useScope: 'global',
  inheritLocale: true,
});
const title = useTitle();
const selectedDashboard = useStore(selectedDashboardStore);

const currentRoute = useRoute();
const routes = [
  {
    name: 'dashboard',
    icon: Dashboard,
    path: '/dashboard',
  },
  {
    name: 'events',
    icon: Events,
    path: '/dashboard/events',
  },
  {
    name: 'integrations',
    icon: Integrations,
    path: '/dashboard/integrations',
  },
  {
    name: 'settings',
    icon: Settings,
    path: '/dashboard/settings',
  },
  {
    name: 'commands',
    icon: Commands,
    path: '/dashboard/commands',
  },
  {
    name: 'timers',
    icon: Timers,
    path: '/dashboard/timers',
  },
  { 
    name: 'moderation',
    icon: Sword,
    path: '/dashboard/moderation',
  },
  {
    name: 'users',
    icon: Users,
    path: '/dashboard/users',
  },
  {
    name: 'keywords',
    icon: Keywords,
    path: '/dashboard/keywords',
  },
  {
    name: 'variables',
    icon: Variables,
    path: '/dashboard/variables',
  },
  {
    name: 'greetings',
    icon: Greetings,
    path: '/dashboard/greetings',
  },
  {
    name: 'overlays',
    icon: Overlays,
    path: '/dashboard/overlays',
  },
  {
    name: 'files',
    path: '/dashboard/files',
  },
  {
    name: 'quotes',
    icon: Quotes,
    path: '/dashboard/quotes',
  },
];

function logOut() {
  localStorage.removeItem('accessToken');
  localStorage.removeItem('refreshToken');
  userStore.set(null);
  router.push('/');
}
</script>

<template>
  <aside class="w-54 hidden sm:block shadow-md min-w-max border-r border-stone-700 h-screen scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600">
    <div class="sidebar-header flex items-center justify-center py-1">
      <div>
        <router-link
          to="/"
          class="mt-5 inline-flex flex-row items-center text-xl"
        >
          <svg
            width="30"
            height="29"
            viewBox="0 0 20 19"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
            class="mr-2"
          >
            <path
              d="M11.22 2.84008H13.8V5.18008C13.8 6.50008 13.98 6.94008 14.36 7.32008C14.72 7.68008 15.32 7.84008 15.86 7.84008H17.06C17.44 7.84008 17.96 7.76008 18.26 7.60008C18.64 7.40008 18.88 7.14008 19.02 6.74008C19.18 6.34008 19.28 5.42008 19.32 4.58008C18.76 4.38008 17.96 4.00008 17.54 3.64008C17.52 4.38008 17.5 5.02008 17.46 5.30008C17.42 5.58008 17.34 5.70008 17.26 5.74008C17.2 5.78008 17.06 5.80008 16.96 5.80008H16.54C16.42 5.80008 16.32 5.78008 16.26 5.72008C16.2 5.64008 16.2 5.46008 16.2 5.10008V0.760078H8.95999V2.82008C8.95999 4.18008 8.73999 5.92008 6.93999 7.18008C7.45999 7.44008 8.49999 8.12008 8.89999 8.52008C10.86 7.04008 11.22 4.72008 11.22 2.88008V2.84008ZM4.81999 0.0800781C3.95999 1.42008 2.15999 3.06008 0.599994 4.04008C0.979994 4.54008 1.51999 5.52008 1.77999 6.04008C3.67999 4.82008 5.73999 2.84008 7.05999 0.980077L4.81999 0.0800781ZM5.43999 4.26008C4.23999 6.24008 2.23999 8.24008 0.399994 9.50008C0.799994 10.0201 1.43999 11.2601 1.65999 11.7801C2.23999 11.3401 2.81999 10.8201 3.41999 10.2401V18.7201H5.71999V7.72008C6.39999 6.88008 7.03999 6.00008 7.55999 5.14008L5.43999 4.26008ZM14.9 10.8001C14.36 11.8201 13.64 12.7201 12.8 13.5201C11.96 12.7201 11.28 11.8201 10.78 10.8001H14.9ZM16.54 8.56008L16.1 8.66008H7.57999V10.8001H10.4L8.61999 11.3401C9.25999 12.7001 10.04 13.8801 11 14.9001C9.63999 15.7601 8.11999 16.4001 6.47999 16.8001C6.91999 17.2801 7.47999 18.1801 7.77999 18.8001C9.59999 18.2601 11.3 17.5001 12.8 16.5001C14.24 17.5401 15.98 18.3201 18.04 18.8001C18.38 18.1601 19.02 17.2001 19.56 16.6801C17.7 16.3401 16.06 15.7401 14.7 14.9401C16.22 13.4801 17.4 11.6201 18.14 9.28008L16.54 8.56008Z"
              fill="white"
            />
          </svg>
          Tsuwari
        </router-link>
      </div>
    </div>
    <div class="px-1 dropend">
      <a
        id="dropdownMenuButton1e"
        href="#"
        class="flex items-center mt-1 text-sm py-4 px-2 h-12 overflow-hidden text-white text-ellipsis whitespace-nowrap rounded border-l-2 border-white transition duration-300 ease-in-out ripple-surface-primary dropdown-toggle"
        data-bs-toggle="dropdown"
        aria-expanded="false"
      >
        <div class="flex items-center">
          <div class="shrink-0">
            <img
              :src="selectedDashboard?.twitch?.profile_image_url ?? user?.profile_image_url"
              class="rounded-full w-10"
              alt="Dashboard Avatar"
            >
          </div>
          <div class="grow ml-3">
            <p class="text-sm font-semibold">{{ selectedDashboard?.twitch?.display_name ?? user?.display_name }}</p>
          </div>
        </div>
      </a>
      <ul
        class="
          dropdown-menu
          min-w-max
          absolute
          hidden
          bg-[#202122]
          text-base
          z-50
          float-left
          py-2
          list-none
          text-left
          text-white
          rounded-lg
          shadow-lg
          mt-1
          hidden
          m-0
          bg-clip-padding
          border-none
          w-48
        "
        style="width: max-content;"
        aria-labelledby="dropdownMenuButton1e"
      >
        <li
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
              hover:bg-[#121212]"
        >
          <div
            class="cursor-pointer"
            @click="setSelectedDashboard(dashboard)"
          >
            <img
              class="w-6 rounded-full inline"
              :src="dashboard?.twitch?.profile_image_url ?? dashboard.twitch?.profile_image_url"
            >
            <span class="ml-4">{{ dashboard.twitch.display_name }}</span>
          </div>
        </li>
        <ul />
      </ul>
    </div>

    <ul class="relative px-1 mt-3">
      <li
        v-for="(route, index) in routes"
        :key="index"
      >
        <RouterLink
          :to="route.path"
          class="flex items-center mt-1 text-sm py-4 px-6 h-12 overflow-hidden text-white text-ellipsis whitespace-nowrap hover:bg-[#202122] rounded border-slate-300 transition duration-300 ease-in-out ripple-surface-primary"
          :class="{
            'bg-neutral-700': currentRoute.path === route.path,
          }"
          @click="title = route.name"
        >
          <span
            v-if="route.icon"
            class="w-3 h-3 mr-3"
          >
            <img :src="route.icon">
          </span>
             
          <span>{{ t(`pages.${route.name}.sidebarName`) }}</span>
        </RouterLink>
      </li>
    </ul>
  </aside>
</template>
