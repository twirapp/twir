<template>
  <section class="bg-black-15 border-b border-b-black-25">
    <div class="container">
      <div class="flex w-[720px] mx-auto">
        <ul class="inline-flex justify-between w-full py-[18px]">
          <li v-for="item in stats" :key="item.id" class="inline-flex flex-col items-center">
            <span class="text-4xl font-medium leading-[130%]">{{ item.stat }}</span>
            <span class="text-gray-60 leading-normal">{{ item.description }}</span>
          </li>
        </ul>
      </div>
    </div>
  </section>
</template>

<script lang="ts" setup>
import { stats } from '@/data/landing/statInfo.js';
</script>
