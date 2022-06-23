<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { HelixStreamData } from '@twurple/api/lib/index.js';
import { useIntervalFn  } from '@vueuse/core';
import { useAxios } from '@vueuse/integrations/useAxios';
import { intervalToDuration, formatDuration } from 'date-fns';
import { ref, watch } from 'vue';

import { api } from '@/plugins/api';
import { localeStore } from '@/stores/locale';
import { selectedDashboardStore } from '@/stores/userStore';

function setLocale(v: string) {
  localeStore.set(v);
}

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
  <nav class="relative w-full flex flex-wrap items-center justify-between py-3 text-white shadow border-b border-stone-700">
    <div class="container-fluid w-full flex flex-wrap items-center justify-between px-6">
      <div
        v-if="stream"
        class="container-fluid flex space-x-2"
      >
        <p>Viewers: <span class="font-bold">{{ stream.viewer_count }}</span></p>
        <p>Category: <span class="font-bold">{{ stream.game_name }}</span></p>
        <p>Title: <span class="font-bold">{{ stream.title }}</span></p>
        <p>Mesages sended: <span class="font-bold">{{ stream.parsedMessages }}</span></p>
        <p>Uptime: <span class="font-bold">{{ uptime }}</span></p>
      </div>
      <div
        v-else
        class="container-fluid flex space-x-2"
      >
        Stream Offline
      </div>
      <div>
        <div class="locale-changer">
          <select 
            v-model="$i18n.locale"
            class="form-control px-3 py-1.5 text-gray-700 rounded select select-sm"
            @change="setLocale($i18n.locale)"
          >
            <option
              v-for="(lang, i) in ['en', 'ru']"
              :key="`Lang${i}`"
              :value="lang"
            >
              {{ lang }}
            </option>
          </select>
        </div>
      </div>
      
      <!--<div class="flex justify-center">
        <div>
          <div class="dropdown relative">
            <button
              id="dropdownMenuButton1"
              class="
          dropdown-toggle
          px-6
          py-2.5
          bg-blue-600
          text-white
          font-medium
          text-xs
          leading-tight
          uppercase
          rounded
          shadow
          hover:bg-blue-700 
          focus:outline-none focus:ring-0
          transition
          duration-150
          ease-in-out
          flex
          items-center
          whitespace-nowrap
        "
              type="button"
              data-bs-toggle="dropdown"
              aria-expanded="false"
            >
              Dropdown button
              <svg
                aria-hidden="true"
                focusable="false"
                data-prefix="fas"
                data-icon="caret-down"
                class="w-2 ml-2"
                role="img"
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 320 512"
              >
                <path
                  fill="currentColor"
                  d="M31.3 192h257.3c17.8 0 26.7 21.5 14.1 34.1L174.1 354.8c-7.8 7.8-20.5 7.8-28.3 0L17.2 226.1C4.6 213.5 13.5 192 31.3 192z"
                />
              </svg>
            </button>
            <ul
              class="
          dropdown-menu
          min-w-max
          absolute
          hidden
          bg-white
          text-base
          z-50
          float-left
          py-2
          list-none
          text-left
          rounded
          shadow
          mt-1
          hidden
          m-0
          bg-clip-padding
          border-none
        "
              aria-labelledby="dropdownMenuButton1"
            >
              <li>
                <a
                  class="
              dropdown-item
              text-sm
              py-2
              px-4
              font-normal
              block
              w-full
              whitespace-nowrap
              bg-transparent
              text-gray-700
              hover:bg-gray-100
            "
                  href="#"
                >Action</a>
              </li>
              <li>
                <a
                  class="
              dropdown-item
              text-sm
              py-2
              px-4
              font-normal
              block
              w-full
              whitespace-nowrap
              bg-transparent
              text-gray-700
              hover:bg-gray-100
            "
                  href="#"
                >Another action</a>
              </li>
              <li>
                <a
                  class="
              dropdown-item
              text-sm
              py-2
              px-4
              font-normal
              block
              w-full
              whitespace-nowrap
              bg-transparent
              text-gray-700
              hover:bg-gray-100
            "
                  href="#"
                >Something else here</a>
              </li>
            </ul>
          </div>
        </div>
      </div>-->
    </div>
  </nav>
</template>