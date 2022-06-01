<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { useQuasar } from 'quasar';
import { computed, ref } from 'vue';
import { useRouter } from 'vue-router';

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

const isMini = ref(false);
const isShow = ref(false);
function toggleMini () {
  isMini.value = !isMini.value;
}

const showDrawer = computed({
  get() {
    return $q.screen.gt.sm ? true : isShow.value;
  },
  set(v: boolean) {
    isShow.value = v;
  },
});

const $q = useQuasar();
const user = useStore(userStore);
const selectedDashboard = useStore(selectedDashboardStore);

const router = useRouter();

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

function logOut() {
  localStorage.removeItem('accessToken');
  localStorage.removeItem('refreshToken');
  userStore.set(null);
  router.push('/');
}
</script>

<template>
  <q-layout view="lHh lpR lFf">
    <q-header
      bordered
      class="bg-dark text-white"
    >
      <q-toolbar>
        <q-btn
          dense
          flat
          round
          icon="menu"
          @click="$q.screen.gt.sm ? toggleMini() : isShow = !isShow"
        />

        <q-toolbar-title>
          <q-avatar>
            <img src="https://cdn.quasar.dev/logo-v2/svg/logo-mono-white.svg">
          </q-avatar>
          Title
        </q-toolbar-title>
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="showDrawer"
      dark
      :mini="isMini && $q.screen.gt.sm"
      side="left"
      bordered
      :width="200"
    >
      <q-toolbar class="text-white">
        <q-toolbar-title>Tsuwari</q-toolbar-title>
      </q-toolbar>
      
      <q-list padding>
        <q-item
          clickable
          to="/dashboard"
          exact
          label="Dashboardd"
        >
          <q-item-section avatar>
            <q-avatar>
              <img :src="selectedDashboard?.twitch?.profile_image_url ?? user?.profile_image_url">
            </q-avatar>
          </q-item-section>
          <q-item-section class="text-white">
            {{ selectedDashboard.twitch.display_name }}
          </q-item-section>

          <q-menu
            fit
            anchor="center middle"
            self="top left"
            max-width="500px"
          >
            <div class="row no-wrap q-pa-md">
              <div class="column">
                <q-list>
                  <q-item
                    v-for="dashboard of user?.dashboards"
                    :key="dashboard.channelId"
                    clickable
                    :class="{'btn-disabled': selectedDashboard.channelId === dashboard.channelId}"
                    @click="setSelectedDashboard(dashboard)"
                  >
                    <q-item-section avatar>
                      <q-avatar
                        color="red"
                        size="36px"
                        text-color="white"
                      >
                        <img
                          :src="dashboard?.twitch?.profile_image_url ?? dashboard.twitch?.profile_image_url"
                        >
                      </q-avatar>
                    </q-item-section>
                    <q-item-section>
                      {{ dashboard.twitch.display_name }}
                    </q-item-section>
                  </q-item>
                </q-list>
              </div>

              <q-separator
                vertical
                inset
                class="q-mx-lg"
              />

              <div class="column items-center">
                <q-avatar size="72px">
                  <img :src="user?.profile_image_url">
                </q-avatar>

                <div class="text-subtitle1 q-mt-md q-mb-xs">
                  {{ user?.display_name ?? user?.login }}
                </div>

                <q-btn
                  v-close-popup
                  color="red"
                  label="Logout"
                  push
                  size="sm"
                  @click="logOut"
                />
              </div>
            </div>
          </q-menu>
        </q-item>
        <q-separator />
        <q-item
          v-for="route of routes"
          :key="route.name"
          clickable
          :to="route.path"  
          exact
          active-class="bg-grey-9 text-white"
        >
          <q-item-section
            v-if="route.icon"
            avatar
          >
            <q-avatar size="24px">
              <img
                :src="route.icon"
              >
            </q-avatar>
          </q-item-section>
          <q-item-section>{{ route.name }}</q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <slot />
    </q-page-container>

    <q-footer
      elevated
      class="bg-grey-8 text-white"
    >
      <q-toolbar>
        <q-toolbar-title>
          <q-avatar>
            <img src="https://cdn.quasar.dev/logo-v2/svg/logo-mono-white.svg">
          </q-avatar>
          <div>Title</div>
        </q-toolbar-title>
      </q-toolbar>
    </q-footer>
  </q-layout>
</template>
