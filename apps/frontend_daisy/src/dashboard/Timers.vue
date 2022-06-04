<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { Timer } from '@tsuwari/prisma';
import { useTitle } from '@vueuse/core';
import type { SetOptional } from 'type-fest';
import { ref } from 'vue';

import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';


const title = useTitle();
title.value = 'Tsuwari - Timers';

type TimerType = SetOptional<Omit<Timer, 'responses' | 'channelId'> & { edit?: boolean, responses: string[] }, 'id'>

const selectedDashboard = useStore(selectedDashboardStore);
// const { execute, data } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/timers`, api, { immediate: false });
const timers = ref<Array<TimerType>>([]);
selectedDashboardStore.subscribe(() => setTimers());

async function setTimers() {
  const { data } = await api(`/v1/channels/${selectedDashboard.value.channelId}/timers`);
  timers.value = data;
}

function insert() {
  timers.value.unshift({
    name: 'My cool timer',
    enabled: true,
    last: 0,
    timeInterval: 60,
    messageInterval: 0,
    responses: [],
    edit: true,
  });
}

async function deleteTimer(index: number, id?: string) {
  if (id) {
    await api.delete(`/v1/channels/${selectedDashboard.value.channelId}/timers/${id}`);
  }

  timers.value = timers.value.filter((_, i) => i !== index);
}

async function saveTimer(timer: TimerType, index: number) {
  let data;
  if (timer.id) {
    const request = await api.put(`/v1/channels/${selectedDashboard.value.channelId}/timers/${timer.id}`, timer);  
    data = request.data;
  } else {
    const request = await api.post(`/v1/channels/${selectedDashboard.value.channelId}/timers`, timer);
    data = request.data;
  }

  timers.value[index] = data;
}
</script>

<template>
  <div class="p-1">
    <div class="flow-root">
      <div class="float-left btn btn-primary btn-sm w-full mb-1 md:w-auto rounded">
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

        <label
          class="ml-1 text-white"
          @click="insert"
        >Add new timer</label>
      </div>

      <input
        type="text"
        placeholder="Search by keyword..."
        class="float-right rounded input input-sm input-bordered w-full md:w-60"
      >
    </div>

    <div class="grid lg:grid-cols-3 md:grid-cols-2 grid-cols-1 gap-2">
      <div
        v-for="timer, timerIndex of timers"
        :key="timerIndex"
        class="card card-compact bg-base-300 drop-shadow-lg rounded"
      >
        <div
          v-if="timer.edit"
          class="card-body"
        >
          <label class="label">
            <span class="label-text">Timer name</span>
            <input
              v-model="timer.enabled"
              type="checkbox"
              class="toggle"
              checked
            >
          </label>
          <input
            v-model="timer.name"
            type="text"
            placeholder="Great timer"
            class="input input-bordered w-full input-sm mb-5 rounded"
          >

          <span class="label-text">Timer time interval (seconds)</span>
          <input
            v-model="timer.timeInterval"
            type="number"
            placeholder="60"
            class="input input-bordered w-full input-sm mb-5 rounded"
          >

          <span class="label-text">Timer messages interval</span>
          <input
            v-model="timer.messageInterval"
            type="number"
            placeholder="0"
            class="input input-bordered w-full input-sm mb-5 rounded"
          >

          <div class="max-h-40 pr-3 scrollbar-thin overflow-auto scrollbar scrollbar-thumb-gray-900 scrollbar-track-gray-600">
            <label
              v-for="response, responseIndex in timer.responses"
              :key="responseIndex"
            >

              <label class="input-group mt-1">
                <input
                  v-model.lazy="timer.responses[responseIndex]"
                  type="text"
                  placeholder="response of timer"
                  class="input input-bordered input-sm w-full rounded"
                >
                <span
                  class="btn btn-error btn-sm rounded no-animation"
                  @click="timer.responses?.splice(responseIndex, 1)"
                >X</span>
              </label>
            </label>
          </div>

          <div class="card-actions flex justify-between">
            <div>
              <button
                class="btn btn-primary btn-sm rounded"
                @click="() => {
                  timer.edit = false
                  if (!timer.id) timers = timers.filter((_, i) => i !== timerIndex)
                }"
              >
                Cancel
              </button>
            </div>
            <div>
              <button
                class="btn btn-success btn-sm w-full mr-2 sm:w-auto rounded"
                @click="timer.responses.push('')"
              >
                Add response
              </button>
              <div class="dropdown dropdown-top dropdown-left">
                <label
                  tabindex="0"
                  class="btn btn-sm btn-error rounded"
                >Delete</label>
                <div
                  tabindex="0"
                  class="dropdown-content card card-compact w-64 p-2 shadow bg-base-300"
                >
                  <h3 class="card-title text-white">
                    Are you sure?
                  </h3>
                  <div class="card-actions">
                    <div class="btn-group w-full">
                      <button
                        class="btn w-1/2 btn-sm rounded"
                        @click="(e) => {
                          (e.target as HTMLElement).blur()
                        }"
                      >
                        No
                      </button>
                      <button
                        class="btn w-1/2 btn-sm rounded"
                        @click="(e) => {
                          (e.target as HTMLElement).blur()
                          deleteTimer(timerIndex, timer.id)
                        }"
                      >
                        Yes
                      </button>
                    </div>
                  </div>
                </div>
              </div>
              <button
                class="btn btn-success btn-sm ml-2 rounded"
                @click="saveTimer(timer, timerIndex)"
              >
                Save
              </button>
            </div>
          </div>
        </div>

        <div
          v-if="!timer.edit"
          class="card-body"
        >
          {{ timer.name }}
        </div>

        <div
          v-if="!timer.edit"
          class="card-actions justify-end m-3"
        >
          <button
            class="btn btn-primary btn-sm rounded"
            @click="timer.edit = true"
          >
            Edit
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
