<!-- eslint-disable vue/no-v-html -->
<template>
  <section class="overflow-hidden">
    <div class="container max-w-[1200px] relative">
      <div class="inline-flex flex-col items-center w-full pt-32 pb-40 min-2xl:pt-40 min-2xl:pb-48">
        <h1 class="gradient-title" v-html="t('sections.firstScreen.title')"></h1>
        <p
          class="
            text-center
            leading-[145%]
            text-gray-70 text-[22px]
            pt-9
            max-w-[600px]
            animate-fadeInDown
            animation-delay-200
            opacity-0
            z-10
          "
        >
          {{ t('tagline') }}
        </p>
        <div
          class="
            relative
            z-10
            inline-grid
            grid-flow-col
            gap-x-4
            mt-12
            opacity-0
            animate-fadeIn
            animation-delay-500
          "
        >
          <a class="learn-more-btn">
            {{ t('buttons.learnMore') }}
          </a>
          <a class="start-for-free-btn">
            {{ t('buttons.startForFree') }}
          </a>
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
          backgroundImage: `url('${WavesSvg}')`,
        }"
      ></div>

      <img
        :src="RhombusSvg"
        class="absolute top-[15%] right-[16%] select-none max-xl:right-[11%] max-lg:hidden"
      />
      <img
        :src="LightingSvg"
        class="absolute top-[41%] left-[14%] select-none max-xl:left-[13%] max-lg:hidden"
      />
      <img
        :src="SmileBotSvg"
        class="absolute right-[6%] select-none bottom-[37%] max-xl:right-[2%] max-lg:hidden"
      />
      <div
        :style="{
          backgroundImage: `url('${PinkBlob}')`,
        }"
        class="blurry-blob -right-[510px] -bottom-[510px] h-[1100px] w-[1050px]"
      />
      <div
        :style="{
          backgroundImage: `url('${BlueBlob}')`,
        }"
        class="blurry-blob -left-[510px] -bottom-[510px] h-[1100px] w-[1075px]"
      />
    </div>
  </section>
</template>

<script lang="ts" setup>
import BlueBlob from '@/assets/blob-blue.png';
import PinkBlob from '@/assets/blob-pink.png';
import LightingSvg from '@/assets/Lighting.svg';
import RhombusSvg from '@/assets/Rhombus.svg';
import SmileBotSvg from '@/assets/SmileBot.svg';
import WavesSvg from '@/assets/Waves.svg';
import useTranslation from '@/hooks/useTranslation.js';

const t = useTranslation<'landing'>();
</script>

<style lang="postcss">
.gradient-title {
  @apply font-bold text-center max-w-[690px] animate-fadeInDown opacity-0 z-10;
  letter-spacing: -2px;
  font-size: 86px;
  line-height: 1.1;

  & > span {
    background: linear-gradient(258.67deg, #d34bf4 4.22%, #905bff 73.01%);

    background-clip: text;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
}

.start-for-free-btn {
  @apply inline-flex
    bg-purple-60 
    px-6
    py-3
    rounded-full
    justify-center
    items-center
    leading-tight
    text-center
    cursor-pointer
    hover:bg-purple-55
    transition-colors
    duration-150;
}

.learn-more-btn {
  @apply inline-flex
    text-white-95
    bg-black-10
    border border-black-25
    px-6
    py-3
    rounded-full
    justify-center
    leading-tight
    text-center
    cursor-pointer
    transition-colors
    hover:bg-black-15;
}

.blurry-blob {
  @apply absolute blur-md animate-fadeIn select-none -z-[5] bg-no-repeat;
}
</style>
