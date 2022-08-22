<template>
  <div class="px-7 py-8 bg-black-15 rounded-[10px] inline-grid gap-y-5">
    <h6 class="text-[32px] font-medium leading-[120%]">
      {{ title }}
    </h6>
    <p class="text-[17px] text-gray-70 leading-normal">
      {{ description }}
    </p>
    <a :href="actionHref" class="text-purple-80 inline-flex items-center">
      {{ actionText }}
      <TswArrowIcon arrowType="narrow" size="24px" class="stroke-purple-80 ml-[7px]" />
    </a>
  </div>
</template>

<script lang="ts" setup>
import { TswArrowIcon } from '@tsuwari/ui-components';

defineProps<{
  title: string;
  description: string;
  actionText: string;
  actionHref: string;
}>();
</script>
