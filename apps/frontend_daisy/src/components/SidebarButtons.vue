<script lang="ts" setup>

import { useStore } from '@nanostores/vue';
import { RouterLink, useRoute, useRouter } from 'vue-router';

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
import Timers from '@/assets/sidebar/timers.svg?url';
import Users from '@/assets/sidebar/users.svg?url';
import Variables from '@/assets/sidebar/variables.svg?url';
import { selectedDashboardStore, userStore, setSelectedDashboard } from '@/stores/userStore';

const currentRoute = useRoute();

const user = useStore(userStore);
const selectedDashboard = useStore(selectedDashboardStore);

const routes = [
  {
    name: 'Dashboard',
    icon: Dashboard,
    path: '/dashboard',
  },
  {
    name: 'Events',
    icon: Events,
    path: '/dashboard/events',
  },
  {
    name: 'Integrations',
    icon: Integrations,
    path: '/dashboard/integrations',
  },
  {
    name: 'Settings',
    icon: Settings,
    path: '/dashboard/settings',
  },
  {
    name: 'Commands',
    icon: Commands,
    path: '/dashboard/commands',
  },
  {
    name: 'Timers',
    icon: Timers,
    path: '/dashboard/timers',
  },
  {
    name: 'Users',
    icon: Users,
    path: '/dashboard/users',
  },
  {
    name: 'Keywords',
    icon: Keywords,
    path: '/dashboard/keywords',
  },
  {
    name: 'Variables',
    icon: Variables,
    path: '/dashboard/variables',
  },
  {
    name: 'Greetings',
    icon: Greetings,
    path: '/dashboard/greetings',
  },
  {
    name: 'Overlays',
    icon: Overlays,
    path: '/dashboard/overlays',
  },
  {
    name: 'Files',
    /* icon: Files, */
    path: '/dashboard/files',
  },
  {
    name: 'Quotes',
    icon: Quotes,
    path: '/dashboard/quotes',
  },
];
</script>

<template>
  <div class="dropdown dropdown-right">
    <!-- <label
      for="my-modal"
      class="btn w-full modal-button rounded mb-4"
      tabindex="0"
    >
      <span class="text-lg"><img
        class="w-4 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2"
        :src="selectedDashboard?.twitch?.profile_image_url ?? user?.profile_image_url"
      ></span>
      <span class="ml-3">{{ selectedDashboard?.twitch?.display_name ?? user?.display_name }}</span>-
    </label>-->
    <ul
      tabindex="0"
      class="dropdown-content rounded menu p-2 shadow bg-base-200 w-52 space-y-1"
    >
      <li
        v-for="dashboard of user?.dashboards"
        :key="dashboard.channelId"
        :class="{'btn-disabled': selectedDashboard.channelId === dashboard.channelId}"
      >
        <span @click="setSelectedDashboard(dashboard)"><img
          class="w-4 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2"
          :src="dashboard?.twitch?.profile_image_url ?? dashboard.twitch?.profile_image_url"
        >{{ dashboard.twitch.display_name }}</span>
      </li>
    </ul>
  </div>
  
  <li
    v-for="(route, index) in routes"
    :key="index"
    class="my-px"
  >
    <RouterLink
      :to="route.path"
      class="h-10 rounded btn btn-sm w-full"
      :class="{
        'btn-ghost': currentRoute.path === route.path,
        'btn-ghost': currentRoute.path !== route.path,
      }"
    >
      <span>{{ route.name }}</span>
    </RouterLink>
  </li>
</template>
