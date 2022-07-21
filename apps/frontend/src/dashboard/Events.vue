<script lang="ts" setup>
import { useTitle } from '@vueuse/core';

import NewEvent from '../assets/buttons/new_event.svg?component';

const title = useTitle();
title.value = 'Tsuwari - Events';
</script>

<template>
  <div>
    <div class="flex">
      <div>
        <div class="border-b border-gray-700 border-r h-screen rounded w-40">
          <h2 class="font-bold mb-2 outline-none pt-2 text-center text-lg">
            Events
          </h2>

          <ul class="menu">
            <!-- @TODO -->
            <li class="border-r">
              <span>Follow</span>
            </li>
            <li><span>Tip</span></li>
          </ul>
        </div>

        <div class="border-b border-gray-700 flex-1 h-screen w-56">
          <div class="mb-5">
            <h2 class="font-bold mb-2 pt-2 text-center text-lg">
              Follow event
            </h2>
            <p class="px-2 text-center text-gray-300">
              Triggering when channel got new follower
            </p>
          </div>

          <div class="bg-base-200 h-4/5 m-1 overflow-y-auto p-1 rounded">
            <!-- @TODO -->
            <div class="my-1">
              <p class="text-xl">
                $username
              </p>
              <p class="text-gray-300">
                Show sender name
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="hidden p-1 sm:block w-full">
        <div class="gap-2 grid grid-cols-1 lg:grid-cols-3 md:grid-cols-2">
          <div class="bg-base-200 card card-compact drop-shadow rounded">
            <div class="card-body">
              <span class="label-text">Action</span>
              <select class="rounded select select-sm w-full">
                <option
                  disabled
                  selected
                >
                  Choose action
                </option>
                <option>Example #1</option>
                <option>Example #2</option>
              </select>

              <span class="label-text">Filter of operion</span>
              <input
                type="text"
                placeholder="$username === 'lwjerri'"
                class="input input-bordered input-sm mb-5 rounded w-full"
              >

              <span class="label-text">Message for sending</span>
              <input
                type="text"
                placeholder="$username, thank u <3"
                class="input input-bordered input-sm rounded w-full"
              >

              <div class="card-actions justify-end">
                <button class="btn btn-error btn-sm rounded">
                  Delete
                </button>
                <button class="btn btn-primary btn-sm rounded">
                  Save
                </button>
              </div>
            </div>
          </div>

          <div
            class="bg-transparent border-2 border-solid card card-compact cursor-pointer hover:border-sky-500 rounded"
          >
            <div class="card-body grid place-items-center">
              <NewEvent />

              <p class="text-lg">
                Add new action
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="md:hidden p-1 w-full">
      <div class="gap-2 grid grid-cols-1 lg:grid-cols-3 md:grid-cols-2">
        <div class="bg-base-200 card card-compact drop-shadow rounded">
          <div class="card-body">
            <span class="label-text">Action</span>
            <select class="select select-sm w-full">
              <option
                disabled
                selected
              >
                Choose action
              </option>
              <option>Example #1</option>
              <option>Example #2</option>
            </select>

            <span class="label-text">Filter of operion</span>
            <input
              type="text"
              placeholder="$username === 'lwjerri'"
              class="input input-bordered input-sm mb-5 w-full"
            >

            <span class="label-text">Message for sending</span>
            <input
              type="text"
              placeholder="$username, thank u <3"
              class="input input-bordered input-sm w-full"
            >

            <div class="card-actions justify-end">
              <button class="btn btn-error btn-sm">
                Delete
              </button>
              <button class="btn btn-primary btn-sm">
                Save
              </button>
            </div>
          </div>
        </div>

        <div class="bg-transparent border-2 border-solid card card-compact cursor-pointer hover:border-sky-500 rounded">
          <div class="card-body grid place-items-center">
            <NewEvent />

            <p class="text-lg">
              Add new action
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
