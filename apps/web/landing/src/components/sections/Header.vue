<template>
  <header class="flex container py-3 items-center justify-between">
    <div class="inline-grid items-center grid-flow-col gap-x-[10px] p-2">
      <img src="@/assets/NewLogo.svg" />
      <span class="font-medium text-xl">Tsuwari</span>
    </div>
    <nav>
      <ul class="inline-grid grid-flow-col gap-x-2">
        <li v-for="item in menuItems" :key="item.id">
          <a :href="item.href" class="leading-tight px-3 py-[10px]">{{ item.name }}</a>
        </li>
      </ul>
    </nav>
    <div class="inline-grid grid-flow-col gap-x-3 items-center">
      <LangSelect />
      <a href="#" class="inline-flex bg-purple-60 px-4 py-[10px] rounded-md leading-tight">Login</a>
    </div>
  </header>
</template>

<script lang="ts" setup>
import LangSelect from '@/components/LangSelect.vue';
import type { NavMenuItem } from '@/types/navMenu';

defineProps<{ menuItems: NavMenuItem[] }>();
</script>
