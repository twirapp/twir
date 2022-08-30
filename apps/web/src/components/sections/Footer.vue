<template>
  <footer class="bg-black-15">
    <div class="container max-w-[1020px]">
      <div class="grid grid-flow-col grid-cols-2">
        <div class="inline-grid gap-y-6 py-10 px-8 justify-start">
          <div class="inline-flex flex-col">
            <div class="inline-flex">
              <img src="@/assets/NewLogo.svg" class="w-[30px] h-[30px]" />
              <span class="text-2xl ml-3 font-medium">Tsuwari</span>
            </div>
            <p class="text-sm text-gray-60 max-w-xs mt-[14px]">
              Created by streamers. Made for streamers. Used by streamers. For streamers with love.
            </p>
          </div>
          <div class="inline-grid grid-flow-col gap-x-8">
            <ul class="inline-grid grid-flow-col gap-x-6">
              <li v-for="item in menuItems" :key="item.id">
                <a :href="rt(item.href)" class="leading-[130%]">{{ rt(item.name) }}</a>
              </li>
            </ul>
            <ul class="inline-grid grid-flow-col gap-x-3">
              <li v-for="item in socials" :key="item.id">
                <a :href="item.href">
                  <TswIcon :name="item.media" size="24px" class="stroke-gray-70" />
                </a>
              </li>
            </ul>
          </div>
          <span class="text-sm text-gray-60">© Tsuwari 2022. All rights reserved.</span>
        </div>
        <div class="inline-flex flex-col py-10 px-8 items-start">
          <h6 class="text-2xl font-medium mb-5">
            Subscribe for updates
          </h6>
          <p class="text-sm text-gray-60 mb-6">
            Non rhoncus, neque arcu, commodo malesuada sed porttitor dictumst integer. Suscipit dictum quam ut blandit amet.
          </p>
          <form class="relative inline-flex max-w-[320px] w-full">
            <input type="email" class=" w-full bg-black-10 text-sm rounded-full px-[18px] py-3 focus:outline-none" placeholder="Type your email" />
            <button type="submit" class="inline-flex rounded-full p-2 bg-purple-60 z-[1] right-[4px] top-[4px] absolute">
              <TswArrowIcon arrowName="ArrowNarrowShort" class="stroke-white-100 w-5 h-5" />
            </button>
          </form>
        </div>
      </div>
    </div>
  </footer>
</template>

<script lang="ts" setup>
import { TswIcon, TswArrowIcon } from '@tsuwari/ui-components';
import { useI18n } from 'vue-i18n';

import type { NavMenuItem } from '@/types/navMenu.js';
import type { SocialMediaItem } from '@/types/socialMedia.js';

defineProps<{ menuItems: NavMenuItem[], socials: SocialMediaItem[] }>();

const { rt } = useI18n();
</script>
