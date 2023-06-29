<!-- eslint-disable vue/no-v-html -->
<template>
  <section class="overflow-hidden">
    <div class="container max-w-[1200px] relative">
      <div class="hero-wrapper">
        <h1 class="gradient-title" v-html="t('sections.firstScreen.title')"></h1>
        <p class="hero-tagline">
          {{ t('tagline') }}
        </p>
        <div class="hero-buttons">
          <TswButton
            :text="t('buttons.learnMore')"
            :isRounded="true"
            href="#"
            size="lg"
            class="justify-center"
            variant="solid-gray"
          />
          <TswButton
            :text="isUserLoggedIn ? 'Go to dashboard' : t('buttons.startForFree')"
            :isRounded="true"
            size="lg"
            class="justify-center"
            variant="solid-purple"
            @click="() => (isUserLoggedIn ? redirectToDashboard() : redirectToLogin())"
          />
        </div>
      </div>
      <div
        class="
          absolute
          -bottom-[438px]
          mx-auto
          left-0
          right-0
          select-none
          w-full
          bg-center
          h-[780px]
          -z-[10]
          bg-no-repeat
        "
        :style="{
          backgroundImage: cssURL(WavesSvg),
        }"
      ></div>
      <RhombusSvg
        class="hero-bg-image top-[15%] right-[16%] max-xl:right-[11%] animation-delay-300"
      />
      <LightingSvg
        class="hero-bg-image top-[41%] left-[14%] max-xl:left-[11%] animation-delay-400"
      />
      <SmileBotSvg
        class="hero-bg-image right-[6%] bottom-[37%] max-xl:right-[3%] animation-delay-600"
      />
      <div
        :style="{
          backgroundImage: cssURL(PinkBlob),
        }"
        class="blurry-blob -right-[510px] -bottom-[510px] h-[1100px] w-[1050px]"
      ></div>
      <div
        :style="{
          backgroundImage: cssURL(BlueBlob),
        }"
        class="blurry-blob -left-[510px] -bottom-[510px] h-[1100px] w-[1075px]"
      ></div>
    </div>
  </section>
</template>

<script lang="ts" setup>
import { cssURL, TswButton } from '@twir/ui-components';
import { isClient } from '@vueuse/core';
import { ref, watch } from 'vue';

import BlueBlob from '@/assets/blob-blue.png';
import PinkBlob from '@/assets/blob-pink.png';
import LightingSvg from '@/assets/Lighting.svg?component';
import RhombusSvg from '@/assets/Rhombus.svg?component';
import SmileBotSvg from '@/assets/SmileBot.svg?component';
import WavesSvg from '@/assets/Waves.svg';
import { redirectToDashboard, redirectToLogin, useUserProfile } from '@/services/auth';
import { useTranslation } from '@/services/locale';

const { t } = useTranslation<'landing'>();

const useUserStatus = () => {
  const isUserLoggedIn = ref<boolean>(false);
  if (!isClient) return isUserLoggedIn;

  const { isFetching, data } = useUserProfile();

  if (!isFetching && data.value !== undefined) {
    isUserLoggedIn.value = true;
    return isUserLoggedIn;
  }
  const stopWatch = watch(
    () => isFetching.value,
    (isFetching) => {
      if (isFetching === false) {
        stopWatch();
        isUserLoggedIn.value = data.value === undefined ? false : true;
      }
    },
  );

  return isUserLoggedIn;
};

const isUserLoggedIn = useUserStatus();
</script>

<style lang="postcss">
.gradient-title {
  @apply leading-[1.1] -tracking-[2px] font-bold text-center
    animate-fadeInDown opacity-0 z-10
    max-w-[690px]
    min-lg:text-[86px] min-sm:text-[74px] min-xs:text-[62px] text-[56px];

  & > span {
    background: linear-gradient(258.67deg, #d34bf4 4.22%, #905bff 73.01%);

    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
}

.hero-wrapper {
  @apply inline-flex
    flex-col
    items-center
    w-full
    pt-32
    pb-40
    min-2xl:pt-36 min-2xl:pb-48
    max-sm:pt-24
    px-10
    max-sm:px-6;

  @media screen and (max-width: 440px) {
    @apply pt-20 pb-32;
  }
}

.hero-tagline {
  @apply text-center
    leading-[145%]
    text-gray-70
    min-sm:text-[22px]
    text-[20px]
    min-md:pt-9
    pt-6
    max-w-[600px]
    animate-fadeInDown
    animation-delay-200
    opacity-0
    z-10;
}

.hero-buttons {
  @apply relative
    z-10
    inline-grid
    grid-flow-col
    gap-x-3
    mt-12
    max-sm:mt-9
    opacity-0
    animate-fadeIn
    animation-delay-500
    max-sm:w-full max-sm:grid-flow-row max-sm:gap-y-3;
}

.blurry-blob {
  @apply absolute blur-md animate-fadeInLong select-none -z-[5] bg-no-repeat;
}

.hero-bg-image {
  @apply absolute
    select-none
    max-lg:hidden
    bg-no-repeat bg-cover
    animate-fadeInLong
    opacity-0;
}
</style>
