<template>
  <Popover
    as="div"
    class="inline-block relative"
  >
    <PopoverButton
      class="relative youtube-player-btn-icon"
      as="button"
    >
      <span class="-right-2 -top-2 absolute bg-[#644EE8] h-5 inline-flex items-center justify-center min-w-[20px] px-[5px] rounded-full text-center text-xs">
        {{ playlistQueue.length > 99 ? '99+' : playlistQueue.length }}
      </span>
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
        class="absolute bg-[#2C2C2C] divide-[#525252] divide-y focus:outline-none mt-2 origin-top-right right-0 ring-1 ring-black ring-opacity-5 rounded shadow-lg tsw-dropdown-menu w-72 z-50"
      >
        <li
          v-for="(video, index) in playlistQueue"
          :key="video.id"
          class="flex gap-x-3 hover:bg-[#383838] items-start px-4 py-3"
        >
          <span class="text-sm">{{ index + 2 }}</span>
          <div class="flex-1 gap-y-2 inline-grid">
            <a
              :href="getYoutubeVideoLink(video.videoId)"
              class="playlist-item-title"
              target="_blank"
            >
              {{ video.title }}
            </a>
            <span class="text-[#AFAFAF] text-xs">Ordered by: {{ video.orderedByName }}</span>
          </div>

          <button @click="() => removeVideo(video.id)">
            <Remove
              class="h-5 hover:cursor-pointer hover:stroke-[#D0D0D0] stroke-[#858585] stroke-[1.5] w-5"
            />
          </button>
        </li>
      </PopoverPanel>
    </Transition>
  </Popover>
</template>

<script lang="ts" setup>
import { Popover, PopoverButton, PopoverPanel } from '@headlessui/vue';

import Playlist from '@/assets/icons/playlist.svg?component';
import Remove from '@/assets/icons/remove.svg?component';
import type { RequestedSong } from '@/functions/useYoutubeSocketPlayer.js';

defineProps<{
  playlistQueue: readonly RequestedSong[]
  removeVideo: (videoId: string) => void;
}>();

const getYoutubeVideoLink = (videoId: string) => `https://youtube.com/watch?v=${videoId}`;
</script>

<style lang="postcss">
.playlist-item-title {
  @apply hover:underline leading-tight text-sm overflow-hidden text-ellipsis;

  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.tsw-dropdown-menu {
  @apply max-h-72 overflow-y-auto;
  ::-webkit-scrollbar-corner {
    background: rgba(0, 0, 0, 0.5);
  }

  scrollbar-width: thin;

  &::-webkit-scrollbar {
    width: 8px;
    height: 12px;
  }

  &::-webkit-scrollbar-track {
    background: rgba(0, 0, 0, 0);
  }

  &::-webkit-scrollbar-thumb {
    background-color: #787878;
    border-radius: 9999px;
    border: 2px solid rgba(0, 0, 0, 0);
    background-clip: padding-box;
  }
}
</style>