<template>
  <div v-if="isError" class="flex h-screen select-none">
    <div class="m-auto space-y-2 text-center text-white">
      <h1 class="font-bold text-6xl">Oops!</h1>
      <p class="text-lg">
        {{ error }}
      </p>
      <button
        class="
          bg-[#9146FF]
          duration-150
          ease-in-out
          focus:outline-none
          focus:ring-0
          hover:bg-[#772CE8]
          inline-block
          leading-tight
          px-10
          py-2.5
          rounded
          shadow
          text-lg text-white
          transition
          uppercase
        "
      >
        Home
      </button>
    </div>
  </div>
  <div v-if="isFetching" class="flex h-screen select-none">
    <div class="m-auto">
      <div class="flex items-center text-white">
        <div class="animate-spin border-4 h-8 inline-block rounded-full w-8" role="status" />
        <span class="ml-2 text-lg">Loading...</span>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { redirectToDashboard, twitchLoginHook } from '@/services/auth.service.js';

const params = new URLSearchParams(window.location.search);

const { error, isError, isFetching, onSuccessLogin } = twitchLoginHook(params);

onSuccessLogin(() => {
  redirectToDashboard();
});
</script>
