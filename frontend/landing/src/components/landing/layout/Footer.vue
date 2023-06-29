<template>
  <footer class="bg-black-15 max-md:bg-black-10 min-md:px-6">
    <div class="container max-w-[1020px]">
      <div
        class="
          min-md:py-8
          w-full
          flex
          max-md:flex-col
          justify-between
          items-center
          min-md:border-b min-md:border-black-25
        "
      >
        <div
          class="
            flex-1 flex
            max-md:w-full
            max-md:py-[18px]
            max-md:justify-center
            max-md:border-b
            max-md:border-b-black-25
          "
        >
          <div class="inline-flex items-center min-md:mr-auto">
            <div class="w-[28px] h-[28px]" :style="{ backgroundImage: cssURL(TsuwariLogo) }" />
            <span class="text-[22px] ml-3 font-medium">Twir</span>
          </div>
        </div>
        <NavMenu
          menuClass="inline-grid grid-flow-col gap-x-6 max-md:py-6 max-sm:grid-flow-row max-sm:justify-items-center max-sm:text-sm max-sm:gap-y-[18px]"
          menuItemClass="leading-[130%]"
        />
        <div class="flex-1 flex max-md:bg-black-15 max-md:w-full max-md:justify-center max-md:py-3">
          <ul class="inline-grid grid-flow-col gap-x-4 min-md:ml-auto">
            <li v-for="item in socials" :key="item.id">
              <a :href="item.href">
                <TswIcon :name="item.type" :height="20" :width="20" class="fill-gray-70" />
              </a>
            </li>
          </ul>
        </div>
      </div>
      <div class="text-xs text-gray-60 w-full text-center py-4">
        {{ t('sections.footer.rights', { year: new Date().getFullYear() }) }}
      </div>
    </div>
  </footer>
</template>

<script lang="ts" setup>
import { TswIcon } from '@twir/ui-components';
import { cssURL } from '@twir/ui-components';

import TsuwariLogo from '@/assets/brand/TsuwariInCircle.svg';
import NavMenu from '@/components/landing/layout/NavMenu.vue';
import { socials } from '@/data/landing/socialMedia.js';
import { useTranslation } from '@/services/locale';

const { t } = useTranslation<'landing'>();
</script>
