<template>
  <div class="inline-flex flex-col overflow-hidden rounded-[10px] bg-black-17">
    <div class="px-5 py-5 inline-flex border-b border-black-25 items-center">
      <div
        class="w-11 h-11 rounded-full mr-[14px] bg-contain"
        :style="{
          backgroundImage: cssURL(avatarUrl),
        }"
      />
      <div class="inline-grid gap-y-[6px] justify-items-start">
        <span
          class="text-lg leading-[130%] whitespace-nowrap overflow-hidden text-ellipsis max-w-full"
        >
          @{{ username }}
        </span>
        <span class="inline-grid grid-flow-col gap-x-[1px]">
          <TswIcon
            v-for="i in 5"
            :key="i"
            name="Star"
            size="18px"
            :class="i > rating ? 'fill-gray-50' : 'fill-purple-70'"
          />
        </span>
      </div>
    </div>
    <div class="px-5 pt-3 pb-5">
      <p class="comment">
        {{ comment }}
      </p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { TswIcon } from '@tsuwari/ui-components';

import { cssURL } from '@/utils/css';

defineProps<{ username: string; rating: number; comment: string; avatarUrl: string }>();
</script>

<style lang="postcss" scoped>
.comment {
  @apply leading-normal text-gray-70;

  max-height: 72px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
