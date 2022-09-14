<template>
  <section class="bg-black-15 py-36">
    <div class="container max-w-[1040px]">
      <div class="flex justify-between items-center">
        <div class="inline-flex flex-col">
          <span class="uppercase tracking-wider text-purple-80 text-lg mb-4">
            {{ t('sections.integrations.preTitle') }}
          </span>
          <h2 class="text-5xl font-semibold max-w-xl leading-[125%]">
            {{ t('sections.integrations.title') }}
          </h2>
          <p class="text-lg text-gray-70 max-w-xl mt-6">
            {{ t('sections.integrations.description') }}
          </p>
        </div>
        <div class="grid grid-cols-3 content-end gap-3" style="direction: rtl">
          <div
            v-for="item in integrationLogos"
            :key="item.__name"
            class="p-6 rounded-xl bg-black-20"
          >
            <component :is="item" width="64px" height="64px" />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script lang="ts" setup>
import { Dota2, Faceit, LastFm, Spotify, VK } from '@tsuwari/ui-icons/logos';

import useTranslation from '@/hooks/useTranslation.js';

const integrationLogos = [Dota2, Faceit, LastFm, Spotify, VK];

const t = useTranslation<'landing'>();
</script>
