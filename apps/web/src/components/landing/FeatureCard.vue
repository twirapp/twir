<template>
  <div ref="card" class="feature-card">
    <h3>{{ title }}</h3>
    <p>{{ description }}</p>
    <a :href="actionHref">
      {{ t('buttons.learnMore') }}
      <TswArrowIcon arrowName="ArrowNarrow" />
    </a>
    <Transition>
      <div
        v-if="!isOutside && isReady"
        class="absolute -z-[1] select-none bg-contain bg-no-repeat"
        :style="{
          backgroundImage: blobBgUrl,
          top: blobTop,
          left: blobLeft,
          height: `${blobHeight / 2}px`,
          width: `${blobWidth / 2}px`,
        }"
      />
    </Transition>
  </div>
</template>

<script lang="ts" setup>
import { TswArrowIcon } from '@tsuwari/ui-components';
import { useImage, useMouseInElement } from '@vueuse/core';
import { computed, ref } from 'vue';

import PinkBlob from '@/assets/blob-pink.png';
import useTranslation from '@/hooks/useTranslation';

defineProps<{
  title: string;
  description: string;
  actionHref: string;
}>();

const t = useTranslation<'landing'>();

const card = ref<HTMLElement | null>(null);

const blobBgUrl = `url('${PinkBlob}')`;

const { elementX, elementY, isOutside } = useMouseInElement(card);

const { state, isReady } = useImage({
  src: PinkBlob,
});

const blobHeight = computed(() => (isReady.value && state.value ? state.value.height : 0));

const blobWidth = computed(() => (isReady.value && state.value ? state.value.width : 0));

const blobTop = computed(() => {
  return `${isOutside.value ? 0 : elementY.value - blobHeight.value / 4}px`;
});

const blobLeft = computed(() => {
  return `${isOutside.value ? 0 : elementX.value - blobWidth.value / 4}px`;
});
</script>

<style lang="postcss" scoped>
.feature-card {
  @apply p-7 rounded-[10px] inline-grid gap-y-6 relative z-10 overflow-hidden border border-black-20 bg-black-15 bg-opacity-30 hover:scale-[1.02] transition-transform duration-300;

  & > h3 {
    @apply text-[32px] font-medium leading-[120%];
  }

  & > p {
    @apply text-[17px] text-gray-70 leading-normal;

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

.v-enter-active,
.v-leave-active {
  transition: opacity 0.3s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
