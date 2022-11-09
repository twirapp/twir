<script lang="ts" setup>
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import VuePlyr from '@skjnldsv/vue-plyr';
import { until } from '@vueuse/core';
import { ref, onMounted } from 'vue';
import '@skjnldsv/vue-plyr/dist/vue-plyr.css';

import Next from '@/assets/icons/next.svg?component';
import Pause from '@/assets/icons/pause.svg?component';
import Play from '@/assets/icons/play.svg?component';
import Refresh from '@/assets/icons/refresh.svg?component';
import Stop from '@/assets/icons/stop.svg?component';
import { NAMESPACES, nameSpaces } from '@/plugins/socket';

const plyr = ref();
const queue = ref<any[]>([]);
const isAlreadyPlayed = ref<boolean>(false);
const paused = ref(true);
const currentTrack = ref<Record<string, any> | undefined>();
const isReady = ref(false);

const socket = nameSpaces.get(NAMESPACES.YOUTUBE)!;
onMounted(() => {
  plyr.value.player.on('ready', () => {
    isReady.value = true;
  });

  plyr.value.player.on('playing', () => {
    paused.value = false;
  });
  plyr.value.player.on('pause', () => {
    paused.value = true;
  });

  plyr.value.player.on('ended', () => {
    isReady.value = false;
    paused.value = false;
    socket.emit('skip', currentTrack.value!.id);
    currentTrack.value = undefined;

    const newTrack = queue.value.shift();
    if (newTrack) {
      setTrack(newTrack);
      play();
    }
  });
  socket.emit('currentQueue', (tracks) => {
    queue.value = tracks;
    console.log(tracks);
    if (queue.value.length) {
      setTrack(queue.value.shift());
    }
  });
  socket.on('newTrack', (track) => {
    const current = currentTrack.value ? true : false;

    if (current) {
      queue.value.push(track);
    } else {
      setTrack(track);
    }

    if (!current && isAlreadyPlayed.value) {
      play();
    }

    console.log('newTrack', queue.value);
  });
  socket.on('removeTrack', (track) => {
    queue.value = queue.value.filter((t: any) => t.videoId !== track.videoId);
  });
});

function stop() {
  if (queue.value.length == 0) {
    plyr.value.player.stop();
    plyr.value.player.source = {
      type: 'video',
      sources: [],
    };
  } else {
    const nextTrack = queue.value.shift();
    if (nextTrack) {
      queue.value.push(nextTrack);
      setTrack(nextTrack);
      plyr.value.player.ready = false;
      play();
    }
  }
}

async function waitPlayer() {
  await until(isReady).toBe(true);
}

async function play() {
  await waitPlayer();
  playHelper();
}

function pause() {
  plyr.value.player.pause();
}

// TODO: CALL SKIP SOMEWHERE
function skip() {
  if (currentTrack.value) {
    socket.emit('skip', currentTrack.value.id);
  }
}

function playHelper() {
  if (!paused.value || !isAlreadyPlayed.value) {
    console.log('restarting...');
    plyr.value.player.restart();
  }
  plyr.value.player.muted = false;
  plyr.value.player.play();
  isAlreadyPlayed.value = true;
}

function setTrack(track) {
  currentTrack.value = track;
  plyr.value.player.source = {
    type: 'video',
    sources: [
      {
        src: track.videoId,
        provider: 'youtube',
      },
    ],
  };
}

function playClick() {
  play();
}
</script>

<template>
  <div class="block break-inside card rounded shadow text-white">
    <h2 class="border-b border-gray-700 card-title flex font-bold justify-center outline-none p-2">
      <p>title</p>
    </h2>
    <div class="p-4 w-full">
      <vue-plyr
        ref="plyr"
        :options="{
          controls: [
            'progress',
            'current-time',
            'mute',
            'volume',
            'captions',
            'settings',
            'pip',
            'airplay',
            'fullscreen',
          ],
          clickToPlay: false,
          hideControls: false,
          keyboard: { focused: false, global: false },
          invertTime: false,
          debug: true,
        }">
        <video controls playsinline />
      </vue-plyr>
      <div class="flex justify-center mt-3">
        <button
          :disabled="!currentTrack"
          :class="{ 'bg-neutral-900': currentTrack }"
          class="border-2 border-zinc-700 p-2 rounded-xl"
          :style="{ cursor: currentTrack ? 'pointer' : 'not-allowed' }"
          @click="stop">
          <Stop />
        </button>
        <button
          :disabled="!currentTrack"
          :class="{ 'bg-neutral-900': currentTrack }"
          class="border-2 border-zinc-700 ml-1 p-2 rounded-xl"
          :style="{ cursor: currentTrack ? 'pointer' : 'not-allowed' }"
          @click="paused ? play() : pause()">
          <Play v-if="paused" />
          <Pause v-else />
        </button>
        <button
          :disabled="queue.length <= 1"
          :class="{ 'bg-neutral-900': queue.length > 1 }"
          :style="{ cursor: queue.length > 1 ? 'pointer' : 'not-allowed' }"
          class="border-2 border-zinc-700 ml-1 p-2 rounded-xl"
          @click="() => {}">
          <Next />
        </button>
      </div>
    </div>
  </div>
</template>

Advertisement
