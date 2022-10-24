<template>
  <div ref="card" class="feature-card">
    <h3>{{ title }}</h3>
    <p>{{ description }}</p>
    <a :href="actionHref">
      {{ t('buttons.learnMore') }}
      <TswIcon name="ArrowNarrow" />
    </a>
    <FeatureCardBgBlob v-if="windowWidth >= 768" :card="card" />
  </div>
</template>

<script lang="ts" setup>
import { TswIcon } from '@tsuwari/ui-components';
import { useWindowSize } from '@vueuse/core';
import { ref } from 'vue';

import FeatureCardBgBlob from '@/components/landing/FeatureCardBgBlob.vue';
import useTranslation from '@/hooks/useTranslation';

defineProps<{
  title: string;
  description: string;
  actionHref: string;
}>();

const t = useTranslation<'landing'>();
const { width: windowWidth } = useWindowSize();

const card = ref<HTMLElement | null>(null);
</script>

<style lang="postcss" scoped>
.feature-card {
  @apply min-md:p-7 py-8 min-md:rounded-[10px] inline-grid gap-y-6 relative z-10 overflow-hidden 
  min-md:border border-b border-black-25 min-md:border-black-20 min-md:bg-black-15 min-md:bg-opacity-30 
  min-md:hover:scale-[1.02] min-md:transition-transform min-md:duration-300;

  & > h3 {
    @apply min-md:text-[32px] text-[30px] font-medium leading-[120%];
  }

  & > p {
    @apply min-md:text-[17px] text-gray-70 leading-normal;

    max-height: 74px;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  & > a {
    @apply text-purple-80 inline-flex items-center transition-colors;

    & > svg {
      @apply stroke-purple-80 ml-[5px] w-[24px] h-[24px] transition-colors;
    }

    &:hover {
      @apply text-[#C4BCF5];

      & > svg {
        @apply stroke-[#C4BCF5];
      }
    }
  }
}
</style>
