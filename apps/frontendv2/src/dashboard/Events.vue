<script lang="ts" setup>
import { useTitle } from '@vueuse/core';

const title = useTitle();
title.value = 'Tsuwari - Events';
</script>

<template>
  <div>
    <div class="flex">
      <div>
        <div class="w-40 h-screen rounded border-r border-b border-gray-700">
          <h2 class="mb-2 pt-2 font-bold text-lg text-center outline-none">
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

        <div class="w-56 h-screen flex-1 border-b border-gray-700">
          <div class="mb-5">
            <h2 class="mb-2 pt-2 text-lg font-bold text-center">
              Follow event
            </h2>
            <p class="px-2 text-center text-gray-300">
              Triggering when channel got new follower
            </p>
          </div>

          <div class="bg-base-200 rounded p-1 m-1 h-4/5 overflow-y-auto">
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

      <div class="w-full p-1 hidden sm:block">
        <div class="grid lg:grid-cols-3 md:grid-cols-2 grid-cols-1 gap-2">
          <div class="card card-compact bg-base-200 drop-shadow-lg rounded">
            <div class="card-body">
              <span class="label-text">Action</span>
              <select class="select select-sm w-full rounded">
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
                class="rounded input input-bordered w-full input-sm mb-5"
              >

              <span class="label-text">Message for sending</span>
              <input
                type="text"
                placeholder="$username, thank u <3"
                class="rounded input input-bordered w-full input-sm"
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
            class="card card-compact bg-transparent cursor-pointer rounded border-solid border-2 hover:border-sky-500"
          >
            <div class="card-body grid place-items-center">
              <svg
                width="48"
                height="48"
                viewBox="0 0 48 48"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
              >
                <path
                  d="M38 6H10C7.79086 6 6 7.79086 6 10V38C6 40.2091 7.79086 42 10 42H38C40.2091 42 42 40.2091 42 38V10C42 7.79086 40.2091 6 38 6Z"
                  stroke="white"
                  stroke-width="3.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M24 16V32"
                  stroke="white"
                  stroke-width="3.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
                <path
                  d="M16 24H32"
                  stroke="white"
                  stroke-width="3.5"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                />
              </svg>

              <p class="text-lg">
                Add new action
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="w-full p-1 md:hidden">
      <div class="grid lg:grid-cols-3 md:grid-cols-2 grid-cols-1 gap-2">
        <div class="card card-compact bg-base-200 drop-shadow-lg rounded">
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
              class="input input-bordered w-full input-sm mb-5"
            >

            <span class="label-text">Message for sending</span>
            <input
              type="text"
              placeholder="$username, thank u <3"
              class="input input-bordered w-full input-sm"
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

        <div class="card card-compact bg-transparent cursor-pointer rounded border-solid border-2 hover:border-sky-500">
          <div class="card-body grid place-items-center">
            <svg
              width="48"
              height="48"
              viewBox="0 0 48 48"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                d="M38 6H10C7.79086 6 6 7.79086 6 10V38C6 40.2091 7.79086 42 10 42H38C40.2091 42 42 40.2091 42 38V10C42 7.79086 40.2091 6 38 6Z"
                stroke="white"
                stroke-width="3.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M24 16V32"
                stroke="white"
                stroke-width="3.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
              <path
                d="M16 24H32"
                stroke="white"
                stroke-width="3.5"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>

            <p class="text-lg">
              Add new action
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
