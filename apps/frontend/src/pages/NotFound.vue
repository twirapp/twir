<script lang="ts" setup>
import { useRouter } from 'vue-router';

import NotFound from '../assets/buttons/back.svg?component';

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
        <NotFound />
      </button>
    </div>
  </div>
</template>