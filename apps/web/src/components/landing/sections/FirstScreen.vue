<!-- eslint-disable vue/no-v-html -->
<template>
  <section class="relative overflow-hidden">
    <ClientOnly>
      <template #default>
        <BlurryBlob
          color="purple"
          class="w-[324px] h-[410px] rotate-[28deg] -right-[100px] -bottom-[30px]"
        />
        <BlurryBlob
          color="blue"
          class="w-[308px] h-[420px] -rotate-[26deg] -left-[100px] -bottom-[80px]"
        />
      </template>
    </ClientOnly>

    <div class="container">
      <div class="inline-flex flex-col items-center w-full pt-32 pb-40">
        <h1
          class="
            font-bold
            text-7xl text-center
            gradient-title
            max-w-[690px]
            animate-fadeInDown
            opacity-0
            z-10
          "
          v-html="t('sections.firstScreen.title')"
        ></h1>
        <p
          class="
            text-center
            leading-[145%]
            text-gray-70 text-[22px]
            pt-7
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
            gap-x-3
            min-w-[320px]
            mt-12
            opacity-0
            animate-fadeIn
            animation-delay-500
          "
        >
          <a
            class="
              inline-flex
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
            "
          >
            {{ t('buttons.learnMore') }}
          </a>
          <a
            class="
              inline-flex
              bg-purple-60
              px-6
              py-3
              rounded-full
              justify-center
              leading-tight
              text-center
              cursor-pointer
            "
          >
            {{ t('buttons.startForFree') }}
          </a>
        </div>
      </div>
      <img
        class="absolute -bottom-[438px] -z-10 mx-auto left-0 right-0 select-none"
        src="@/assets/Waves.svg"
      />
      <img src="@/assets/Rhombus.svg" class="absolute top-[15%] right-[24%] select-none" />
      <img class="absolute top-[41%] left-[22%] select-none" src="@/assets/Lighting.svg" />
      <img src="@/assets/SmileBot.svg" class="absolute right-[216px] select-none bottom-[200px]" />
    </div>
  </section>
</template>

<script lang="ts" setup>
import BlurryBlob from '@/components/BlurryBlob.vue';
import ClientOnly from '@/components/ClientOnly.vue';
import { useTranslation } from '@/utils/locales.js';

const t = useTranslation<'landing'>();
</script>

<style lang="postcss">
.gradient-title span {
  background: linear-gradient(258.67deg, #d34bf4 4.22%, #905bff 73.01%);

  background-clip: text;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  line-height: 1.15;
}

@keyframes fadeIn {
  0% {
    opacity: 0;
  }
  100% {
    opacity: 1;
  }
}
</style>
