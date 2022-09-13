<template>
  <div class="px-6 inline-flex justify-items-center items-center w-full">
    <div
      :class="`w-[80px] h-[80px] inline-block bg-contain`"
      :style="{
        backgroundImage: `url('${UserImage}')`,
      }"
    />
    <div class="inline-flex flex-col items-start w-full flex-1 ml-6">
      <div class="inline-flex flex-col justify-items-center">
        <div class="inline-flex items-center w-full mb-2">
          <span class="text-[21px] font-medium leading-tight">
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
            {{ t('sections.team.founder') }}
          </span>
        </div>
        <p class="text-gray-70 max-w-[180px]">
          {{ role }}
        </p>
      </div>
      <ul class="inline-grid grid-flow-col gap-x-3 mt-[16px]">
        <li v-for="(media, index) in socials" :key="index" class="inline-flex">
          <a :href="media.link" class="inline-block" target="_blank" rel="noopener noreferrer">
            <TswIcon :name="media.type" size="22px" class="stroke-gray-70" />
          </a>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { TswIcon } from '@tsuwari/ui-components';

import UserImage from '@/assets/User.svg';
import type { TeamMemberMedia } from '@/data/landing/team.js';
import useTranslation from '@/hooks/useTranslation';

defineProps<{
  name: string;
  role: string;
  socials: TeamMemberMedia[];
  isFounder?: boolean;
}>();

const t = useTranslation<'landing'>();
</script>