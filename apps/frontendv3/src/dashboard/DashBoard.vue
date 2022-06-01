<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { useTimeoutPoll, get, useTitle  } from '@vueuse/core';
import { ref } from 'vue';

import { socketEmit } from '@/plugins/socket';
import { selectedDashboardStore } from '@/stores/userStore';


const title = useTitle();
title.value = 'Tsuwari - Widgets';

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
    console.log('isBotMod', data);
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
  <div class="q-pa-md">
    <div class="row">
      <q-card
        flat
        bordered
        class="col-4"
      >
        <q-card-section>
          <div class="text-h6">
            Status
          </div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-banner
            rounded
            class="bg-primary text-white"
            :class="{ 'bg-deep-orange-6': !isBotMod, 'bg-green-6': isBotMod }"
          >
            <span v-if="isBotMod">
              Bot is mod on the channel
            </span>
            <span v-else>
              We detect bot is not moderator on the channel. Please, mod a bot, or some of functionality may work incorrectly.
            </span>
          </q-banner>
        </q-card-section>

        <q-separator inset />

        <q-card-section class="q-pa-md q-gutter-sm">
          <q-btn
            color="deep-orange"
            label="Leave"
            @click="leaveChannel"
          />
       
          <q-btn
            color="secondary"
            label="Join"
            @click="joinChannel"
          />
        </q-card-section>
      </q-card>
    </div>
  </div>
</template>

<style scoped>
.card {
  background: linear-gradient(0deg, rgba(255, 255, 255, 0.05), rgba(255, 255, 255, 0.05)), #121212;
}
</style>