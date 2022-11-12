<template>
  <div class="block break-inside card overflow-hidden rounded shadow text-white">
    <div class="border-b border-gray-700 flex items-center justify-between px-4 py-3">
      <p class="font-bold ml-1">
        Youtube song request
      </p>
      <div class="gap-x-1 grid-flow-col inline-grid">
        <Popover
          v-if="queueLength > 1"
          as="div"
          class="inline-block relative"
        >
          <PopoverButton
            class="youtube-player-btn-icon"
            as="button"
          >
            <Playlist class="stroke-icon" />
          </PopoverButton>
          <Transition
            enter-active-class="transition duration-100 ease-out"
            enter-from-class="transform scale-95 opacity-0"
            enter-to-class="transform scale-100 opacity-100"
            leave-active-class="transition duration-75 ease-in"
            leave-from-class="transform scale-100 opacity-100"
            leave-to-class="transform scale-95 opacity-0"
          >
            <PopoverPanel
              as="ul"
              class="absolute bg-[#2C2C2C] focus:outline-none mt-2 origin-top-right right-0 ring-1 ring-black ring-opacity-5 rounded shadow-lg w-72 z-50"
              :class="{
                'tsw-dropdown-menu divide-[#525252] divide-y': queueLength > 1
              }"
            >
              <li
                v-for="video, index in queue"
                :key="video.id"
                class="flex gap-x-3 hover:bg-[#383838] items-start px-4 py-3"
              >
                <span class="text-sm">{{ index + 2 }}</span>
                <div class="flex-1 gap-y-2 inline-grid">
                  <a
                    :href="`https://youtube.com/watch?v=${video.videoId}`"
                    class="hover:underline leading-tight text-sm two-lines-max"
                    target="_blank"
                  >{{ video.title }}</a>
                  <span class="text-[#AFAFAF] text-xs">Ordered by: {{ video.orderedByName }}</span>
                </div>
                
                <button
                  @click="() => removeVideo(video.id)"
                >
                  <Remove class="h-5 hover:cursor-pointer hover:stroke-[#D0D0D0] stroke-[#858585] stroke-[1.5] w-5" />
                </button>
              </li>
            </PopoverPanel>
          </Transition>
        </Popover>
        <button class="youtube-player-btn-icon">
          <DotsHorizontal class="fill-icon" />
        </button>
      </div>
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
      class="aspect-video bg-[#2C2C2C] flex flex-col h-full items-center justify-center text-[#AFAFAF] w-full"
    >
      <div
        v-if="isLoadingQueue"
        class="inline-flex items-center"
      >
        <div
          class="animate-spin border-2 h-6 inline-block mr-3 rounded-full spinner-border w-6"
          role="status"
        />
        <span>Loading...</span>
      </div>
      <template v-else>
        <OffVideo class="h-20 mb-2 stroke-[#AFAFAF] stroke-[2] w-20" />
        <span>There is no videos in queue</span>
      </template>
    </div>
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
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';
import { computed } from 'vue';

import VuePlyr from './VuePlyr.vue';


import DotsHorizontal from '@/assets/icons/dots-horizontal.svg?component';
import Next from '@/assets/icons/next.svg?component';
import OffVideo from '@/assets/icons/off-video.svg?component';
import Pause from '@/assets/icons/pause.svg?component';
import Play from '@/assets/icons/play.svg?component';
import Playlist from '@/assets/icons/playlist.svg?component';
import Remove from '@/assets/icons/remove.svg?component';
import { useYoutubeSocketPlayer } from '@/functions/useYoutubeSocketPlayer.js';

const {
  initQueue,
  currentVideo,
  isPaused,
  isActive,
  skipCurrentVideo,
  removeVideo,
  queueWithoutFirst: queue,
  queueLength,
  isLoadingQueue,
} = useYoutubeSocketPlayer();

const isAcitveStyle = computed(() => (isActive.value ? 'block' : 'none'));
</script>

<style lang="postcss">
.plyr {
  display: v-bind(isAcitveStyle);
}

.youtube-player-btn-icon {
  @apply hover:bg-[#2C2C2C] inline-flex p-[6px] hover:cursor-pointer rounded;

  &[aria-expanded='true'] {
    @apply bg-[#2C2C2C];
    & > svg.fill-icon {
      @apply fill-[#D0D0D0];
    }
  }

  &[aria-expanded='true'] {
    @apply bg-[#2C2C2C];
    & > svg.stroke-icon {
      @apply stroke-[#D0D0D0];
    }
  }

  & > svg {
    @apply h-5 w-5;
  }

  & > svg.fill-icon {
    @apply fill-[#AFAFAF];
  }

  & > svg.stroke-icon {
    @apply stroke-[#AFAFAF];
  }

  &:hover {
    & > svg.fill-icon {
      @apply hover:fill-[#D0D0D0];
    }

     & > svg.stroke-icon {
      @apply hover:stroke-[#D0D0D0];
    }
  }
}
</style>

<style scoped lang="postcss">

.two-lines-max {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.tsw-dropdown-menu {
  @apply max-h-72 overflow-y-auto;
  ::-webkit-scrollbar-corner { background: rgba(0,0,0,0.5); }

  scrollbar-width: thin;
  /* scrollbar-color: var(--scroll-bar-color) var(--scroll-bar-bg-color); */

  &::-webkit-scrollbar {
    width: 8px;
    height: 12px;
  }

  &::-webkit-scrollbar-track {
    background: rgba(0,0,0,0.0);
  }

  &::-webkit-scrollbar-thumb {
    background-color: #787878;
    border-radius: 9999px;
    border: 2px solid rgba(0, 0, 0, 0);
    background-clip: padding-box;
    /* border: 3px solid rgba(0,0,0,0.0); */
  }
}
</style>
