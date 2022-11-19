<template>
  <div class="block break-inside card overflow-hidden rounded shadow text-white">
    <div class="border-b border-gray-700 flex items-center justify-between px-4 py-3">
      <p class="font-bold ml-1">
        Youtube song request
      </p>
      <div class="gap-x-1 grid-flow-col inline-grid">
        <YoutubePlaylistQueue
          v-if="queueLength > 1"
          :playlistQueue="playlistQueue"
          :removeVideo="removeVideo"
        />
        <YoutubeSettingsModal />
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
        <p class="font-medium two-lines-max">
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
import YoutubePlaylistQueue from './YoutubePlaylistQueue.vue';
import YoutubeSettingsModal from './YoutubeSettingsModal.vue';

import Next from '@/assets/icons/next.svg?component';
import OffVideo from '@/assets/icons/off-video.svg?component';
import Pause from '@/assets/icons/pause.svg?component';
import Play from '@/assets/icons/play.svg?component';
import { useYoutubeSocketPlayer } from '@/functions/useYoutubeSocketPlayer.js';

const {
  initQueue,
  currentVideo,
  isPaused,
  isActive,
  skipCurrentVideo,
  removeVideo,
  queueWithoutFirst: playlistQueue,
  queueLength,
  isLoadingQueue,
} = useYoutubeSocketPlayer();

const isActiveStyle = computed(() => (isActive.value ? 'block' : 'none'));
</script>

<style lang="postcss">
.plyr {
  display: v-bind(isActiveStyle);
}

.youtube-player-btn-icon {
  @apply hover:bg-[#2C2C2C] inline-flex p-[6px] hover:cursor-pointer rounded;

  &[aria-expanded='true'] {
    @apply bg-[#2C2C2C];

    & > svg.fill-icon {
      @apply fill-[#D0D0D0];
    }

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
      @apply fill-[#D0D0D0];
    }

    & > svg.stroke-icon {
      @apply stroke-[#D0D0D0];
    }
  }
}
</style>
