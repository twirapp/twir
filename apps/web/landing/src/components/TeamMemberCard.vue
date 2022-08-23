<template>
  <div class="px-6 inline-flex justify-items-center items-center w-full">
    <div class="w-[80px] h-[80px] inline-block bg-[url('../assets/User.svg')] bg-contain" />
    <div class="inline-flex flex-col items-start w-full flex-1 ml-6">
      <div class="inline-grid gap-y-[6px] justify-items-center">
        <div class="inline-flex items-center w-full">
          <span class="text-[21px] font-medium leading-normal">
            {{ name }}
          </span>
          <span
            v-if="isFounder"
            class="
              px-[10px]
              py-[3px]
              rounded-full
              text-yellow-65
              bg-yellow-65 bg-opacity-[0.15]
              text-sm
              ml-2
            "
          >
            Founder
          </span>
        </div>
        <p class="text-gray-70 max-w-[180px]">
          {{ role }}
        </p>
      </div>
      <ul class="inline-grid grid-flow-col gap-x-3 mt-[16px]">
        <li v-for="(social, index) in socials" :key="index" class="inline-flex">
          <a :href="social.link" class="inline-block" target="_blank" rel="noopener noreferrer">
            <TswIcon :name="social.media" size="22px" class="stroke-gray-70" />
          </a>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { TswIcon } from '@tsuwari/ui-components';

import type { TeamMemberMediaLink } from '@/types/teamMember';

withDefaults(
  defineProps<{
    name: string;
    role: string;
    isFounder?: boolean;
    socials: TeamMemberMediaLink[];
  }>(),
  {
    isFounder: false,
  },
);
</script>
