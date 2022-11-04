<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import VuePlyr from '@skjnldsv/vue-plyr';
import { onMounted, onUnmounted, ref } from 'vue';
import '@skjnldsv/vue-plyr/dist/vue-plyr.css';

import { youtubeAutoPlay } from './YoutubePlayerStore';

import Next from '@/assets/icons/next.svg?component';
import Pause from '@/assets/icons/pause.svg?component';
import Play from '@/assets/icons/play.svg?component';
import Refresh from '@/assets/icons/refresh.svg?component';
import Stop from '@/assets/icons/stop.svg?component';
import { nameSpaces, NAMESPACES } from '@/plugins/socket';

const autoPlay = useStore(youtubeAutoPlay);
const plyr = ref();
const queue = ref<any[]>([]);
const paused = ref(true);

const youtube = nameSpaces.get(NAMESPACES.YOUTUBE)!;
onMounted(() => {
  youtube.on('currentQueue', (d) => {
    queue.value = d;
    if (d.length) {
      setTrack(d[0]);
    }
  });

  youtube.on('newTrack', (track) => {
    if (!queue.value.length) {
      setTrack(track);
    }
    queue.value.push(track);
    if (autoPlay.value && queue.value.length === 1) {
      play(true);
    }
  });

  plyr.value.player.on('ended', () => next());
});

function setTrack(track) {
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

function next() {
  const currentTrack = queue.value[0];
  const nextTrack = queue.value[1];
  if (currentTrack) {
    queue.value = queue.value.filter((v) => v.id != currentTrack.id);
    youtube.emit('skip', currentTrack.id);
  }
  if (nextTrack) {
    setTrack(nextTrack);
    play(true);
  }
}

function playHelper() {
  if (!paused.value) {
    plyr.value.player.currentTime = 0;
  }
  plyr.value.player.muted = false;
  plyr.value.player.play();
  paused.value = false;
}

function play(wait = false) {
  if (wait) {
    plyr.value.player.off('ready');
    plyr.value.player.on('ready', () => {
      playHelper();
    });
  } else {
    playHelper();
  }
}

function pause() {
  plyr.value.player.pause();
  paused.value = true;
}

function stop() {
  youtube.emit('skip', queue.value[0].id);
  if (queue.value.length == 1) {
    plyr.value.player.stop();
    plyr.value.player.source = {
      type: 'video',
      sources: [],
    };
    queue.value = [];
    paused.value = false;
  } else {
    pause();
    next();
  }
}
</script>

<template>
  <div class="block break-inside card rounded shadow text-white">
    <h2 class="border-b border-gray-700 card-title flex font-bold justify-center outline-none p-2">
      <p>title</p>
    </h2>
    <div class="p-4 w-full">
      <TransitionGroup name="slide-up">
        <div :key="`player`" v-show="queue.length > 0">
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
            }">
            <video controls playsinline />
          </vue-plyr>
        </div>

        <div
          :key="`spinner`"
          v-show="queue.length === 0"
          class="bg-neutral-900 flex flex-col items-center p-11 show">
          <span><Refresh class="animate-spin" /></span>
          <span>Waiting song requests...</span>
        </div>
      </TransitionGroup>

      <div class="flex justify-center mt-3">
        <button
          :disabled="!queue.length"
          :class="{ 'bg-neutral-900': queue.length }"
          class="border-2 border-zinc-700 p-2 rounded-xl"
          @click="stop">
          <Stop />
        </button>
        <button
          :disabled="!queue.length"
          :class="{ 'bg-neutral-900': queue.length }"
          class="border-2 border-zinc-700 ml-1 p-2 rounded-xl"
          @click="paused ? play() : pause()">
          <Play v-if="paused && queue.length" />
          <Pause v-else />
        </button>
        <button
          :disabled="queue.length <= 1"
          :class="{ 'bg-neutral-900': queue.length > 1 }"
          :style="{ cursor: queue.length > 1 ? 'pointer' : 'not-allowed' }"
          class="border-2 border-zinc-700 ml-1 p-2 rounded-xl"
          @click="next">
          <Next />
        </button>
      </div>

      <div v-if="queue.length" class="flex flex-col mt-2">
        <span
          ><span>{{ queue[0].title }}</span></span
        >
        <span
          >Requested by
          <span
            ><b>{{ queue[0].orderedByName }}</b></span
          ></span
        >
      </div>
    </div>
  </div>
</template>

<style>
.plyr {
  max-height: 300px;
}

.show {
}

.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.25s ease-out;
}

.slide-up-enter-from {
  opacity: 0;
  transform: translateY(30px);
}

.slide-up-leave-to {
  opacity: 0;
  transform: translateY(-30px);
}
</style>
