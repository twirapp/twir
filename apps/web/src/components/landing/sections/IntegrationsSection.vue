<template>
  <section
    class="
      relative
      min-lg:py-32 min-md:py-24 min-xl:px-12 min-lg:px-16 min-lg:bg-black-15
      pb-24
      min-md:px-10 min-sm:px-8 min-sm:pt-48
      pt-60
      px-5
    "
  >
    <div class="container max-w-[1040px] relative">
      <div class="flex justify-between items-center relative z-0">
        <div class="inline-flex flex-col z-[1]">
          <span class="uppercase tracking-wider text-purple-80 text-lg mb-4">
            {{ t('sections.integrations.preTitle') }}
          </span>
          <h2
            class="
              min-xl:text-[48px] min-lg:text-[40px]
              tracking-tight
              text-[36px]
              font-semibold
              min-lg:max-w-xl
              max-w-lg
              leading-[125%]
            "
          >
            {{ t('sections.integrations.title') }}
          </h2>
          <p class="text-lg text-gray-70 min-lg:max-w-xl max-w-lg mt-6">
            {{ t('sections.integrations.description') }}
          </p>
        </div>
        <div
          class="
            grid grid-cols-3
            content-end
            gap-3
            rotate-[5deg]
            flex-shrink-0
            max-lg:absolute max-xl:-z-[1]
            min-md:right-0 min-lg:top-0 min-md:-top-[80px] min-sm:-top-[170px]
            -right-[54px]
            min-lg:scale-100
            scale-90
            -top-[240px]
          "
          style="direction: rtl"
        >
          <div
            v-for="item in integrationLogos"
            :key="item.__name"
            class="
              p-[26px]
              rounded-xl
              min-lg:bg-black-20
              bg-black-20 bg-opacity-60
              hover:bg-opacity-0
              border-2 border-opacity-0 border-black-25
              hover:border-opacity-100
              hover:shadow-xl
              transition-all
              duration-100
              max-xl:p-6
            "
          >
            <component :is="item" class="flex h-16 w-16" />
          </div>
        </div>
      </div>
    </div>
    <div
      :style="{ backgroundImage: cssURL(BlueBlob) }"
      class="
        absolute
        -top-[180px]
        max-md:-left-[380px] max-md:-top-[100px]
        -right-[340px]
        w-[760px]
        h-[990px]
        bg-no-repeat bg-contain
        -z-[10]
        opacity-0
        animate-fadeIn
      "
    ></div>
  </section>
</template>

<script lang="ts" setup>
import { Dota2, Faceit, LastFm, Spotify, VK } from '@tsuwari/ui-icons/logos';

import BlueBlob from '@/assets/blob-blue.png';
import useTranslation from '@/hooks/useTranslation.js';
import { cssURL } from '@/utils/css.js';

const integrationLogos = [Dota2, Faceit, LastFm, Spotify, VK];

const t = useTranslation<'landing'>();
</script>
