<script lang="ts" setup>
import { useIntervalFn  } from '@vueuse/core';
import { intervalToDuration, formatDuration } from 'date-fns';
import { ref } from 'vue';

const uptime = ref('');

useIntervalFn(() => {
  uptime.value = formatDuration(intervalToDuration({ start: new Date('2022-05-24T21:56:15Z'), end: Date.now() }));
}, 1000, { immediate: true });

</script>

<template>
  <nav class="relative w-full flex flex-wrap items-center justify-between py-3 text-white shadow-lg border-b border-stone-700">
    <div class="container-fluid w-full flex flex-wrap items-center justify-between px-6">
      <div class="container-fluid flex space-x-2">
        <p>Online: <span class="font-bold">{{ uptime }}</span></p>
        <p>Viewers: <span class="font-bold">158</span></p>
      </div>
    </div>
  </nav>
</template>

<style scoped>
@media (min-width: 1280px) {
  #show-button {
    display: none;
  }
}
</style>
