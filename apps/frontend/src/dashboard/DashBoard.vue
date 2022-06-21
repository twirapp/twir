<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { useTimeoutPoll, get, useTitle  } from '@vueuse/core';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import { socketEmit } from '@/plugins/socket';
import { selectedDashboardStore } from '@/stores/userStore';

const title = useTitle();
title.value = 'Dashboard';

const { t } = useI18n({
  useScope: 'global',
  inheritLocale: true,
});

const isBotMod = ref(false);

const dashboard = useStore(selectedDashboardStore);

selectedDashboardStore.subscribe(() => isBotMod.value = false);

useTimeoutPoll(() => {
  const dash = get(dashboard);
  if (!dash) return;

  socketEmit('isBotMod', {
    channelId: dash.channelId,
    channelName: dash.twitch.login,
    userId: dash.userId,
  }, (data) => {
    isBotMod.value = data.value;
  });
}, 1000, { immediate: true });

function leaveChannel() {
  console.log('leaving');
  socketEmit('botPart', { 
    channelName: selectedDashboardStore.get().twitch.login,
    channelId: selectedDashboardStore.get().channelId,
  });
}

function joinChannel() {
  console.log('joining');
  socketEmit('botJoin', { 
    channelName: selectedDashboardStore.get().twitch.login,
    channelId: selectedDashboardStore.get().channelId,
  });
}

/* watch(isWindowFocused, (v) => {
  if (v) checkIsMod.resume();
  else checkIsMod.pause();
}); */

</script>

<template>
  <div class="m3">
    <div class="p-1">
      <div class="grid lg:grid-cols-3 grid-cols-1 gap-2">
        <div
          class="block rounded-lg card text-white shadow-lg max-w-sm"
        >
          <h2 class="card-title p-2 flex justify-between border-b border-gray-700 outline-none">
            <p>{{ t('pages.dashboard.widgets.status.title') }}</p>

          <!-- <a
              href="/"
              class="btn btn-outline btn-error btn-sm rounded"
            >Remove</a> -->
          </h2>
          <div class="p-4">
            <div
              class="rounded-lg py-5 px-6 text-base mb-4"
              :class="{ 'text-yellow-700 bg-yellow-100': !isBotMod, 'bg-green-100 text-green-700': isBotMod }"
            >
              <div v-if="!isBotMod">
                <div class="text-sm">
                  {{ t('pages.dashboard.widgets.status.notMod' ) }}
                </div>
              </div>
              <div v-else>
                {{ t('pages.dashboard.widgets.status.mod' ) }}
              </div>
            </div>
            <button
              type="button"
              class="inline-block px-6 py-2.5 bg-red-600 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-red-700 hover:shadow-lg focus:bg-red-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-red-800 active:shadow-lg transition duration-150 ease-in-out"
              @click="leaveChannel"
            >
              {{ t('pages.dashboard.widgets.status.buttons.leave') }}
            </button>
            <button
              type="button"
              class="inline-block ml-2 px-6 py-2.5 bg-green-500 text-white font-medium text-xs leading-tight uppercase rounded shadow-md hover:bg-green-600 hover:shadow-lg focus:bg-green-600 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-green-700 active:shadow-lg transition duration-150 ease-in-out"
              @click="joinChannel"
            >
              {{ t('pages.dashboard.widgets.status.buttons.join') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>