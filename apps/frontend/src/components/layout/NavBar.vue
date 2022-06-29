<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { HelixStreamData } from '@twurple/api';
import { useIntervalFn  } from '@vueuse/core';
import { useAxios } from '@vueuse/integrations/useAxios';
import { intervalToDuration, formatDuration } from 'date-fns';
import { ref, watch } from 'vue';

import LanguageSelector from '../LanguageSelector.vue';
import Notification from '../Notification.vue';
import Profile from '../Profile.vue';

import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';

const selectedDashboard = useStore(selectedDashboardStore);
const { execute, data: axiosData } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/streams`, api, { immediate: false });

const stream = ref<HelixStreamData & { parsedMessages: number } | null>(null);
const uptime = ref('');

watch(axiosData, (v) => {
  stream.value = v;
});

selectedDashboardStore.subscribe(async (v) => {
  execute(`/v1/channels/${v.channelId}/streams`);
});

useIntervalFn(() => {
  execute(`/v1/channels/${selectedDashboard.value.channelId}/streams`);
}, 15000);

useIntervalFn(() => {
  if (stream.value) {
    uptime.value = formatDuration(intervalToDuration({ start: new Date(stream.value.started_at), end: Date.now() }));
  }
}, 1000, { immediate: true });
</script>

<template>
  <nav class="border-b border-stone-700 flex flex-wrap items-center justify-between py-3 relative shadow text-white w-full">
    <div class="flex flex-wrap items-center justify-between px-6 w-full">
      <div
        v-if="stream"
        class="flex space-x-2"
      >
        <p>Viewers: <span class="font-bold">{{ stream.viewer_count }}</span></p>
        <p class="hidden lg:block">
          Category: <span class="font-bold">{{ stream.game_name }}</span>
        </p>
        <p class="hidden md:block">
          Title: <span class="font-bold">{{ stream.title.length >= 20 ? stream.title.slice(0, 20) + "..." : stream.title }}</span>
        </p>
      </div>
      <div
        v-else
        class="flex space-x-2"
      >
        Stream Offline
      </div>

      
      <div class="flex flex-row space-x-3.5">
        <LanguageSelector />
      
        <Notification />
        <Profile />
      </div>
    </div>
  </nav>
</template>