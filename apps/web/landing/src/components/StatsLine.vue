<template>
  <div class="bg-black-15 border-b border-b-black-25">
    <div class="container">
      <div class="flex w-[720px] mx-auto">
        <ul class="inline-flex justify-between w-full py-[18px]">
          <li v-for="item in statsInfo" :key="item.id" class="inline-flex flex-col items-center">
            <span class="text-4xl font-medium leading-[130%]">{{ item.stats }}</span>
            <span class="text-gray-60 leading-normal">{{ item.description }}</span>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
const statsInfo: { stats: string; description: string; id: number }[] = [
  { id: 1, stats: '1053+', description: 'Streamers' },
  { id: 2, stats: '51251', description: 'Reviews' },
  { id: 3, stats: '50%', description: 'Aliquam nulla' },
  { id: 4, stats: '64531', description: 'Commads created' },
];
</script>
