<template>
  <div class="inline-flex flex-col overflow-hidden rounded-[10px] bg-black-17 max-w-sm">
    <div class="px-5 py-5 inline-flex border-b border-black-25">
      <img :src="avatarUrl" class="w-11 h-11 rounded-full mr-[14px]" />
      <div>
        <span class="text-lg leading-[130%]">@{{ username }}</span>
        <span></span>
      </div>
    </div>
    <div class="px-5 pt-3 pb-5">
      <p class="leading-normal text-gray-70">
        {{ comment }}
      </p>
    </div>
  </div>
</template>

<script lang="ts" setup>
defineProps<{username: string, rating: number, comment: string, avatarUrl: string}>();
</script>