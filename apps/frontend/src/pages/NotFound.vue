<script lang="ts" setup>
import { useRouter } from 'vue-router';

const router = useRouter();
</script>

<template>
  <div class="flex h-full select-none text-white">
    <div class="m-auto space-y-2 text-center">
      <h1 class="font-bold text-6xl">
        Oops!
      </h1>
      <p
        class="text-lg"
      >
        Page {{ router.currentRoute.value.fullPath === '/' ? '' : router.currentRoute.value.fullPath }} not found.
      </p>

      <button
        class="bg-[#9146FF] duration-150 ease-in-out focus:outline-none focus:ring-0 hover:bg-[#772CE8] inline-block leading-tight px-10 py-2.5 rounded shadow text-lg text-white transition uppercase"
        @click="$router.go(-1)"
      >
        <svg
          class="h-6 w-6"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
          xmlns="http://www.w3.org/2000/svg"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="1"
          d="M12.066 11.2a1 1 0 000 1.6l5.334 4A1 1 0 0019 16V8a1 1 0 00-1.6-.8l-5.333 4zM4.066 11.2a1 1 0 000 1.6l5.334 4A1 1 0 0011 16V8a1 1 0 00-1.6-.8l-5.334 4z"
        /></svg>
      </button>
    </div>
  </div>
</template>