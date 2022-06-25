<script lang="ts" setup>
/* eslint-disable vue/no-v-html */
import { useStore } from '@nanostores/vue';
import { useTimeoutPoll, get, useTitle  } from '@vueuse/core';
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';

import Soon from '@/components/Soon.vue';
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
  <div class="m-1.5 md:m-3">
    <div class="p-1">
      <div class="grid lg:grid-cols-3 grid-cols-1 gap-2">
        <div
          class="block rounded card text-white shadow"
        >
          <h2 class="card-title font-bold p-2 flex justify-center border-b border-gray-700 outline-none">
            <p>{{ t('pages.dashboard.widgets.status.title') }}</p>

          <!-- <a
              href="/"
              class="btn btn-outline btn-error btn-sm rounded"
            >Remove</a> -->
          </h2>
          <div class="p-4 w-full">
            <div
              class="rounded py-5 px-6 text-base mb-4"
              :class="{ 'bg-[#ED4245]': !isBotMod, 'bg-green-600': isBotMod }"
            >
              <div v-if="!isBotMod">
                <div v-html="t('pages.dashboard.widgets.status.notMod' )" />
              </div>
              <div v-else>
                {{ t('pages.dashboard.widgets.status.mod' ) }}
              </div>
            </div>
         
            <div class="flex flex-col space-y-1 md:flex-row md:space-x-1 md:justify-end md:text-right md:space-y-0">
              <button
                type="button"
                class="inline-block px-6 py-2.5 bg-red-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-red-700  focus:outline-none focus:ring-0   transition duration-150 ease-in-out"
                @click="leaveChannel"
              >
                {{ t('pages.dashboard.widgets.status.buttons.leave') }}
              </button>
              <button
                type="button"
                class="inline-block px-6 py-2.5 bg-green-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-green-700  focus:outline-none focus:ring-0   transition duration-150 ease-in-out"
                @click="joinChannel"
              >
                {{ t('pages.dashboard.widgets.status.buttons.join') }}
              </button>
            </div>
          </div>
        </div>

        <div
          class="block rounded card text-white shadow"
        >
          <h2 class="card-title font-bold p-2 flex justify-center border-b border-gray-700 outline-none">
            <p>{{ t('pages.dashboard.widgets.feedback.title') }}</p>
          </h2>
          <div class="p-4 inline-block w-full">
            <Soon :button="false" />
            <!-- <div class="flex justify-center">
              <div class="mb-3 w-full">
                <textarea
                  id="exampleFormControlTextarea1"
                  class="
        form-control
        block
        w-full
        px-3
        py-1.5
        text-base
        font-normal
        text-gray-700
        bg-white bg-clip-padding
        border border-solid border-gray-300
        rounded
        transition
        ease-in-out
        m-0
        focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none
      "
                  rows="3"
                
                  :placeholder="t('pages.dashboard.widgets.feedback.placeholder')"
                />
              </div>
            </div>

            <div class="text-right">
              <button
                type="button"
                class="opacity-60 pointer-events-none inline-block ml-2 px-6 py-2.5 bg-green-600 text-white font-medium text-xs leading-tight uppercase rounded shadow hover:bg-green-700  focus:outline-none focus:ring-0   transition duration-150 ease-in-out"
                
                @click="joinChannel"
              >
                {{ t('pages.dashboard.widgets.feedback.buttons.send') }}
              </button>
            </div> -->
          </div>
        </div>
      </div>
    </div>
  </div>
</template>