<script lang="ts" setup>
import { useStore } from '@nanostores/vue';
import { Greeting } from '@tsuwari/prisma';
import { useTitle } from '@vueuse/core';
import { useAxios } from '@vueuse/integrations/useAxios';
import type { SetOptional } from 'type-fest';
import { ref } from 'vue';

import { api } from '@/plugins/api';
import { selectedDashboardStore } from '@/stores/userStore';


const title = useTitle();
title.value = 'Tsuwari - Greetings';

type GreeTingType = SetOptional<Omit<Greeting, 'channelId'> & { username: string, edit?: boolean }, 'id'>

const selectedDashboard = useStore(selectedDashboardStore);

// const { execute, data } = useAxios(`/v1/channels/${selectedDashboard.value.channelId}/greetings`, api, { immediate: false });
const greetings = ref<Array<GreeTingType>>([]);
selectedDashboardStore.subscribe((v) => setGreetings());

async function setGreetings() {
  const { data } = await api(`/v1/channels/${selectedDashboard.value.channelId}/greetings`);
  greetings.value = data;
}


function insert() {
  greetings.value.unshift({
    username: '',
    userId: '',
    text: '',
    edit: true,
    enabled: true,
  });
}

async function deleteGreeting(index: number, id?: string) {
  if (id) {
    await api.delete(`/v1/channels/${selectedDashboard.value.channelId}/greetings/${id}`);
  }

  greetings.value = greetings.value.filter((_, i) => i !== index);
}

async function saveGreeting(greeting: GreeTingType, index: number) {
  let data;
  if (greeting.id) {
    const request = await api.put(`/v1/channels/${selectedDashboard.value.channelId}/greetings/${greeting.id}`, greeting);  
    data = request.data;
  } else {
    const request = await api.post(`/v1/channels/${selectedDashboard.value.channelId}/greetings`, greeting);
    data = request.data;
  }

  greetings.value[index] = data;
}
</script>

<template>
  <div class="p-1">
    <div class="flow-root">
      <div class="float-left rounded btn btn-primary btn-sm w-full mb-1 md:w-auto">
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
        >Add new greeting</label>
      </div>

      <input
        type="text"
        placeholder="Search by keyword..."
        class="float-right rounded input input-sm input-bordered w-full md:w-60"
      >
    </div>

    <div class="grid lg:grid-cols-3 md:grid-cols-2 grid-cols-1 gap-2">
      <div
        v-for="greeting, greetingIndex of greetings"
        :key="greeting.username"
        class="card card-compact bg-base-200 drop-shadow-lg rounded"
      >
        <div
          v-if="greeting.edit"
          class="card-body"
        >
          <label class="label">
            <span class="label-text">Username</span>
            
            <button
              v-if="greeting.enabled"
              class="label-text-alt rounded btn-outline btn btn-error btn-sm"
              @click="greeting.enabled = !greeting.enabled"
            >Disable</button>
            <button
              v-if="!greeting.enabled"
              class="label-text-alt btn-outline btn btn-success btn-sm rounded"
              @click="greeting.enabled = !greeting.enabled"
            >Enable</button>
          </label>
          <input
            v-model.lazy="greeting.username"
            type="text"
            placeholder="tsuwaribot"
            class="input input-bordered w-full input-sm mb-5 rounded"
          >

          <span class="label-text">Message for sending</span>
          <input
            v-model.lazy="greeting.text"
            type="text"
            placeholder="$(sender), hello!"
            class="input input-bordered w-full input-sm rounded"
          >

          <div class="card-actions flex justify-between">
            <div>
              <button
                class="btn btn-primary btn-sm rounded"
                @click="() => {
                  greeting.edit = false
                  if (!greeting.id) greetings = greetings.filter((_, i) => i !== greetingIndex)
                }"
              >
                Cancel
              </button>
            </div>
            <div>
              <div class="dropdown dropdown-top dropdown-left">
                <label
                  tabindex="0"
                  class="btn btn-sm btn-error rounded"
                >Delete</label>
                <div
                  tabindex="0"
                  class="dropdown-content rounded card card-compact w-64 p-2 shadow bg-base-300"
                >
                  <h3 class="card-title text-white">
                    Are you sure?
                  </h3>
                  <div class="card-actions">
                    <div class="btn-group w-full">
                      <button
                        class="btn w-1/2 rounded btn-sm"
                        @click="(e) => {
                          (e.target as HTMLElement).blur()
                        }"
                      >
                        No
                      </button>
                      <button
                        class="btn w-1/2 rounded btn-sm"
                        @click="(e) => {
                          (e.target as HTMLElement).blur()
                          deleteGreeting(greetingIndex, greeting.id)
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
                @click="saveGreeting(greeting, greetingIndex)"
              >
                Save
              </button>
            </div>
          </div>
        </div>

        <div
          v-if="!greeting.edit"
          class="card-body"
        >
          <div class="card-title">
            {{ greeting.username }}
          </div>
          <div>{{ greeting.text }}</div>
        </div>

        <div
          v-if="!greeting.edit"
          class="card-actions justify-end m-3"
        >
          <button
            class="btn btn-primary btn-sm"
            @click="greeting.edit = true"
          >
            Edit
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
