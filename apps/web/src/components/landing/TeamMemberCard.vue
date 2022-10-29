<template>
  <div class="min-lg:px-6 inline-flex min-lg:justify-items-center min-lg:items-center w-full">
    <UserImage
      class="
        min-sm:w-[80px] min-sm:h-[80px]
        w-[68px]
        h-[68px]
        inline-block
        bg-contain bg-no-repeat
        p-1
      "
    />
    <div class="inline-flex flex-col items-start w-full flex-1 min-md:ml-6 ml-4">
      <div class="inline-flex flex-col justify-items-center">
        <div class="inline-flex items-center w-full mb-2">
          <span class="text-[20px] font-medium leading-tight">
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
        <p class="text-gray-70 min-lg:max-w-[180px] text-sm">
          {{ role }}
        </p>
      </div>
      <ul class="inline-grid grid-flow-col gap-x-[14px] mt-[16px]">
        <li v-for="(media, index) in socials" :key="index" class="inline-flex">
          <a :href="media.link" class="inline-flex" target="_blank" rel="noopener noreferrer">
            <TswIcon
              :name="media.type"
              :height="20"
              :width="20"
              class="stroke-gray-70 fill-gray-70 hover:stroke-white-95 hover:fill-white-95"
            />
          </a>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { TswIcon } from '@tsuwari/ui-components';

import UserImage from '@/assets/User.svg?component';
import type { TeamMemberMedia } from '@/data/landing/team.js';
import { useTranslation } from '@/services/locale/hooks.js';

defineProps<{
  name: string;
  role: string;
  socials: TeamMemberMedia[];
  isFounder?: boolean;
}>();

const { t } = useTranslation<'landing'>();
</script>
