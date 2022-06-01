<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { useTimeoutPoll, get, useWindowFocus, useTitle  } from '@vueuse/core';
import { ref, watch } from 'vue';

import { socketEmit } from '@/plugins/socket';
import { selectedDashboardStore } from '@/stores/userStore';


const title = useTitle();
title.value = 'Tsuwari - Widgets';

const isBotMod = ref(false);

const dashboard = useStore(selectedDashboardStore);

selectedDashboardStore.subscribe(() => isBotMod.value = false);

useTimeoutPoll(() => {
  const dash = get(dashboard);
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
  socketEmit('botPart', { 
    channelName: selectedDashboardStore.get().twitch.login,
    channelId: selectedDashboardStore.get().channelId,
  });
}

function joinChannel() {
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
  <div class="p-1">
    <div class="dropdown mb-5">
      <!-- svelte-ignore a11y-label-has-associated-control -->

      <div
        class="btn btn-primary btn-sm rounded"
        tabindex="0"
      >
        <label>
          <svg
            width="20"
            height="20"
            viewBox="0 0 20 20"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              d="M10 4.16663V15.8333"
              stroke="white"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
            <path
              d="M4.16663 10H15.8333"
              stroke="white"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </label>

        <label class="ml-1 text-white">Add new widget</label>
      </div>

      <ul
        tabindex="0"
        class="dropdown-content menu p-2 bg-base-200 rounded w-52"
      >
        <li><span>Chat</span></li>
        <li><span>Stream</span></li>
        <li><span>Eventlist</span></li>
      </ul>
    </div>

    <div class="grid lg:grid-cols-3 grid-cols-1 gap-2">
      <div class="card rounded card-compact bg-base-200 drop-shadow-lg min-h-96">
        <div class="card-body h-3/4">
          <h2 class="card-title pb-1 border-b border-gray-700 outline-none">
            <p>Status</p>

            <a
              href="/"
              class="btn btn-outline btn-error btn-sm rounded"
            >Remove</a>
          </h2>

          <div>
            <div
              class="alert rounded drop-shadow-lg"
              :class="{ 'alert-warning': !isBotMod, 'alert-success': isBotMod }"
            >
              <div v-if="!isBotMod">
                <svg
                  class="w-16 h-16"
                  fill="none"
                  stroke="currentColor"
                  viewBox="2 0 24 24"
                  xmlns="http://www.w3.org/2000/svg"
                ><path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636"
                /></svg>
                <div class="text-sm">
                  We detect bot is not moderator on the channel. Please, mod a bot, or some of functionality may work incorrectly.
                </div>
              </div>
              <div v-else>
                Bot is mod on the channel
              </div>
            </div>
          </div>
        </div>
        <div class="card-actions justify-end mb-5 mr-5">
          <button
            class="btn btn-primary btn-outline btn-sm rounded"
            @click="leaveChannel"
          >
            Leave
          </button>

          <button
            class="btn btn-primary btn-outline btn-sm rounded"
            @click="joinChannel"
          >
            Join
          </button>
        </div>
      </div>

      <!--<div class="card rounded card-compact bg-base-200 drop-shadow-lg min-h-96">
        <div class="card-body h-3/4">
          <h2 class="card-title pb-1 border-b border-gray-700 outline-none">
            <div>
              <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path
                  d="M16.5 8.58336C16.5029 9.68325 16.2459 10.7683 15.75 11.75C15.162 12.9265 14.2581 13.916 13.1395 14.6078C12.021 15.2995 10.7319 15.6662 9.41667 15.6667C8.31678 15.6696 7.23176 15.4126 6.25 14.9167L1.5 16.5L3.08333 11.75C2.58744 10.7683 2.33047 9.68325 2.33333 8.58336C2.33384 7.26815 2.70051 5.97907 3.39227 4.86048C4.08402 3.7419 5.07355 2.838 6.25 2.25002C7.23176 1.75413 8.31678 1.49716 9.41667 1.50002H9.83333C11.5703 1.59585 13.2109 2.32899 14.4409 3.55907C15.671 4.78915 16.4042 6.42973 16.5 8.16669V8.58336Z"
                  stroke="#6C43B5"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </div>
            <p>Chat</p>

            <a href="/" class="btn btn-outline btn-error btn-sm rounded">Remove</a>
          </h2>
          <p>Twitch chat</p>

          <div class="card-actions pt-1 border-t border-gray-700 justify-end">
            <input type="text" placeholder="Send your message" class="input-sm bg-base-300 input w-full" />
            <button class="btn btn-primary btn-sm">Send</button>
          </div>
        </div>
      </div>-->

      <div class="card rounded card-compact bg-base-200 drop-shadow-lg min-h-96">
        <div class="card-body">
          <h2 class="card-title pb-1 border-b border-gray-700 outline-none">
            <div>
              <svg
                width="20"
                height="20"
                viewBox="0 0 20 20"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M13.3333 9.1665V5.83317M17.5 1.6665H2.5V14.9998H6.66667V18.3332L10 14.9998H14.1667L17.5 11.6665V1.6665ZM9.16667 9.1665V5.83317V9.1665Z"
                  stroke="#6C43B5"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </div>
            <p>Stream</p>

            <a
              href="/"
              class="btn btn-outline btn-error btn-sm rounded"
            >Remove</a>
          </h2>
          <p>Stream</p>
        </div>
      </div>

      <div class="card rounded card-compact bg-base-200 drop-shadow-lg min-h-96">
        <div class="card-body">
          <h2 class="card-title pb-1 border-b border-gray-700 outline-none">
            <div>
              <svg
                width="20"
                height="19"
                viewBox="0 0 20 19"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M10 1.6665L12.575 6.88317L18.3334 7.72484L14.1667 11.7832L15.15 17.5165L10 14.8082L4.85002 17.5165L5.83335 11.7832L1.66669 7.72484L7.42502 6.88317L10 1.6665Z"
                  stroke="white"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>
            </div>
            <p>Eventlist</p>

            <a
              href="/"
              class="btn btn-outline btn-error btn-sm rounded"
            >Remove</a>
          </h2>

          <div
            style="height: 50vh"
            class="overflow-y-auto scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600"
          >
            <table class="table w-full">
              <tbody>
                <tr>
                  <td>
                    LWJerri
                    <br>
                    <span class="badge badge-ghost badge-sm">Follow. Yesterday, 19:30</span>
                  </td>
                  <td>
                    <svg
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M17.3667 3.84172C16.941 3.41589 16.4357 3.0781 15.8795 2.84763C15.3232 2.61716 14.7271 2.49854 14.125 2.49854C13.5229 2.49854 12.9268 2.61716 12.3705 2.84763C11.8143 3.0781 11.309 3.41589 10.8833 3.84172L10 4.72506L9.11666 3.84172C8.25692 2.98198 7.09086 2.49898 5.875 2.49898C4.65914 2.49898 3.49307 2.98198 2.63333 3.84172C1.77359 4.70147 1.29059 5.86753 1.29059 7.08339C1.29059 8.29925 1.77359 9.46531 2.63333 10.3251L3.51666 11.2084L10 17.6917L16.4833 11.2084L17.3667 10.3251C17.7925 9.89943 18.1303 9.39407 18.3608 8.83785C18.5912 8.28164 18.7099 7.68546 18.7099 7.08339C18.7099 6.48132 18.5912 5.88514 18.3608 5.32893C18.1303 4.77271 17.7925 4.26735 17.3667 3.84172V3.84172Z"
                        fill="white"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </td>
                </tr>

                <tr>
                  <td>
                    LWJerri
                    <br>
                    <span class="badge badge-ghost badge-sm">Follow. Yesterday, 19:30</span>
                  </td>
                  <td>
                    <svg
                      width="12"
                      height="20"
                      viewBox="0 0 12 20"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M6 0.833496V19.1668"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M10.1667 4.1665H3.91667C3.14312 4.1665 2.40125 4.47379 1.85427 5.02078C1.30729 5.56776 1 6.30962 1 7.08317C1 7.85672 1.30729 8.59858 1.85427 9.14557C2.40125 9.69255 3.14312 9.99984 3.91667 9.99984H8.08333C8.85688 9.99984 9.59875 10.3071 10.1457 10.8541C10.6927 11.4011 11 12.143 11 12.9165C11 13.6901 10.6927 14.4319 10.1457 14.9789C9.59875 15.5259 8.85688 15.8332 8.08333 15.8332H1"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </td>
                </tr>

                <tr>
                  <td>
                    LWJerri
                    <br>
                    <span class="badge badge-ghost badge-sm">Follow. Yesterday, 19:30</span>
                  </td>
                  <td>
                    <svg
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M10.8333 1.6665L2.5 11.6665H10L9.16667 18.3332L17.5 8.33317H10L10.8333 1.6665Z"
                        fill="white"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </td>
                </tr>

                <tr>
                  <td>
                    LWJerri
                    <br>
                    <span class="badge badge-ghost badge-sm">Follow. Yesterday, 19:30</span>
                  </td>
                  <td>
                    <svg
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M10 18.3332C14.6024 18.3332 18.3334 14.6022 18.3334 9.99984C18.3334 5.39746 14.6024 1.6665 10 1.6665C5.39765 1.6665 1.66669 5.39746 1.66669 9.99984C1.66669 14.6022 5.39765 18.3332 10 18.3332Z"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M18.3333 10H15"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M5.00002 10H1.66669"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M10 4.99984V1.6665"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M10 18.3333V15"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </td>
                </tr>
                <tr>
                  <td>
                    LWJerri
                    <br>
                    <span class="badge badge-ghost badge-sm">Follow. Yesterday, 19:30</span>
                  </td>
                  <td>
                    <svg
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M10 18.3332C14.6024 18.3332 18.3334 14.6022 18.3334 9.99984C18.3334 5.39746 14.6024 1.6665 10 1.6665C5.39765 1.6665 1.66669 5.39746 1.66669 9.99984C1.66669 14.6022 5.39765 18.3332 10 18.3332Z"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M18.3333 10H15"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M5.00002 10H1.66669"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M10 4.99984V1.6665"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M10 18.3333V15"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </td>
                </tr>
                <tr>
                  <td>
                    LWJerri
                    <br>
                    <span class="badge badge-ghost badge-sm">Follow. Yesterday, 19:30</span>
                  </td>
                  <td>
                    <svg
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M10 18.3332C14.6024 18.3332 18.3334 14.6022 18.3334 9.99984C18.3334 5.39746 14.6024 1.6665 10 1.6665C5.39765 1.6665 1.66669 5.39746 1.66669 9.99984C1.66669 14.6022 5.39765 18.3332 10 18.3332Z"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M18.3333 10H15"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M5.00002 10H1.66669"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M10 4.99984V1.6665"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M10 18.3333V15"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </td>
                </tr>
                <tr>
                  <td>
                    LWJerri
                    <br>
                    <span class="badge badge-ghost badge-sm">Follow. Yesterday, 19:30</span>
                  </td>
                  <td>
                    <svg
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M10 18.3332C14.6024 18.3332 18.3334 14.6022 18.3334 9.99984C18.3334 5.39746 14.6024 1.6665 10 1.6665C5.39765 1.6665 1.66669 5.39746 1.66669 9.99984C1.66669 14.6022 5.39765 18.3332 10 18.3332Z"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M18.3333 10H15"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M5.00002 10H1.66669"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M10 4.99984V1.6665"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M10 18.3333V15"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </td>
                </tr>
                <tr>
                  <td>
                    LWJerri
                    <br>
                    <span class="badge badge-ghost badge-sm">Follow. Yesterday, 19:30</span>
                  </td>
                  <td>
                    <svg
                      width="20"
                      height="20"
                      viewBox="0 0 20 20"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M10 18.3332C14.6024 18.3332 18.3334 14.6022 18.3334 9.99984C18.3334 5.39746 14.6024 1.6665 10 1.6665C5.39765 1.6665 1.66669 5.39746 1.66669 9.99984C1.66669 14.6022 5.39765 18.3332 10 18.3332Z"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M18.3333 10H15"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M5.00002 10H1.66669"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M10 4.99984V1.6665"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                      <path
                        d="M10 18.3333V15"
                        stroke="white"
                        stroke-width="2"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
