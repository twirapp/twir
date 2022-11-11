<template>
  <div class="block break-inside card overflow-hidden rounded shadow text-white">
    <div class="border-b border-gray-700 px-5 py-3">
      <p class="font-bold">
        Youtube song request
      </p>
    </div>
    <VuePlyr
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
        ratio: '16:9',
        clickToPlay: false,
        hideControls: true,
        keyboard: { focused: false, global: false },
        invertTime: false,
        debug: false,
      }"
      :style="{
        '--plyr-color-main': '#644EE8',
      }"
      @init="initQueue"
    />
    <div
      v-if="!isActive"
      class="aspect-video bg-[#2C2C2C] flex h-full items-center justify-center text-[#AFAFAF] w-full"
    >
      <span>There is no videos in queue</span>
    </div>
    <!-- <ul>
      <li
        v-for="video in queue"
        :key="video.id"
      >
        <span>{{ video.title }} : {{ video.orderedByName }}</span>
        <button
          class="bg-red-500 text-white"
          @click="() => removeVideo(video.id)"
        >
          Remove
        </button>
      </li>
    </ul>  -->
    <div
      v-if="isActive"
      class="border-[#403D3A] border-t flex items-start p-5"
    >
      <div class="flex-1 gap-y-2 inline-grid mr-5">
        <p class="font-medium">
          {{ currentVideo!.title }}
        </p>
        <span class="text-[#AFAFAF] text-xs">Ordered by: {{ currentVideo!.orderedByName }}</span>
      </div>
      
      <div class="gap-x-3 grid-flow-col inline-grid">
        <component
          :is="isPaused ? Play : Pause"
          class="h-5 hover:cursor-pointer hover:stroke-[#D0D0D0] stroke-[#AFAFAF] w-5"
          @click="isPaused = !isPaused"
        />
        <Next
          class="h-5 hover:cursor-pointer hover:stroke-[#D0D0D0] stroke-[#AFAFAF] w-5"
          @click="skipCurrentVideo"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

import VuePlyr from './VuePlyr.vue';

import Next from '@/assets/icons/next.svg?component';
import Pause from '@/assets/icons/pause.svg?component';
import Play from '@/assets/icons/play.svg?component';
import { useYoutubeSocketPlayer } from '@/functions/useYoutubeSocketPlayer.js';

const { initQueue, currentVideo, isPaused, isActive, skipCurrentVideo, removeVideo, queueWithoutFirst: queue } =
  useYoutubeSocketPlayer();

const isAcitveStyle = computed(() => (isActive.value ? 'block' : 'none'));
</script>

<style>
.plyr {
  display: v-bind(isAcitveStyle);
}
</style>
