<template>
  <div v-if="isLoading" class="flex h-screen select-none">
    <div class="m-auto flex items-center">
      <TswLoader size="lg" />
      <span class="ml-2 text-lg">Loading...</span>
    </div>
  </div>
  <div v-if="isError" class="flex h-screen select-none">
    <div class="m-auto flex flex-col items-center space-y-5 text-center text-white">
      <h1 class="font-bold text-7xl">Oops!</h1>
      <p class="text-lg text-gray-70 inline-block">
        {{ error }}
      </p>
      <div class="inline-grid grid-flow-col gap-x-2">
        <TswButton text="Retry" variant="solid-gray" rightIcon="Reload" @click="redirectToLogin" />
        <TswButton text="Back to home" href="/" rightIcon="Home" />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { TswButton, TswLoader } from '@twir/ui-components';

import { useTwitchAuth, redirectToLogin } from '@/services/auth';

const { isLoading, isError, error } = useTwitchAuth();
</script>
