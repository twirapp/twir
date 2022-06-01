<script lang="ts" setup>
import { useIntervalFn  } from '@vueuse/core';
import { intervalToDuration, formatDuration } from 'date-fns';
import { ref } from 'vue';

import SidebarButtons from './SidebarButtons.vue';

const uptime = ref('');

useIntervalFn(() => {
  uptime.value = formatDuration(intervalToDuration({ start: new Date('2022-05-24T21:56:15Z'), end: Date.now() }));
}, 1000, { immediate: true });

</script>

<template>
  <header class="header border-b border-gray-700">
    <div class="text-white header-content px-5 py-2 flex justify-between items-center flex-row">
      <!-- @TODO -->
      <p class="ml-3">
        Viewers: <span class="font-bold text-base">536</span>
      </p>
      <p class="ml-3">
        Views: <span class="font-bold text-base">1 256 256</span>
      </p>
      <p class="ml-3 hidden lg:block">
        Uptime: <span class="font-bold text-base">{{ uptime }}</span>
      </p>
      <p class="ml-3 hidden sm:block">
        Category: <span class="font-bold text-base">Just Chatting</span>
      </p>
      <p class="ml-3 hidden md:block">
        Title:
        <span class="font-bold text-base">ОБЗОР НА CALL OF DUTY MODERN WARFRAME...</span>
      </p>

      <label
        id="show-button"
        for="profile-popout"
        class="btn btn-ghost btn-circle btn-sm"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-5 w-5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 6h16M4 12h16M4 18h7"
          />
        </svg>
      </label>
    </div>

    <input
      id="profile-popout"
      type="checkbox"
      class="modal-toggle"
    >
    <label
      for="profile-popout"
      class="modal modal-bottom cursor-pointer"
    >
      <label
        class="modal-box relative"
        for=""
      >
        <ul class="flex flex-col w-full">
          <SidebarButtons />

          <div class="divider" />

          <!-- @TODO -->
          <a
            href="/"
            class="btn btn-error btn-sm"
          >Log out</a>
        </ul>
      </label>
    </label>
  </header>
</template>

<style scoped>
@media (min-width: 1280px) {
  #show-button {
    display: none;
  }
}
</style>
