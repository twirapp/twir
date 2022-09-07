<template>
  <footer class="bg-black-15 px-6">
    <div class="container max-w-[1020px]">
      <div class="py-8 w-full flex justify-between items-center border-b border-black-25">
        <div class="flex-1 flex">
          <div class="inline-flex items-center mr-auto">
            <img src="@/assets/NewLogo.svg" class="w-[28px] h-[28px]" />
            <span class="text-[22px] ml-3 font-medium">Tsuwari</span>
          </div>
        </div>
        <NavMenu
          :menuItems="menuItems"
          menuClass="inline-grid grid-flow-col gap-x-6"
          menuItemClass="leading-[130%]"
        />
        <div class="flex-1 flex">
          <ul class="inline-grid grid-flow-col gap-x-3 ml-auto">
            <li v-for="item in socials" :key="item.id">
              <a :href="item.href">
                <TswIcon :name="item.media" size="24px" class="stroke-gray-70" :strokeWidth="1.4" />
              </a>
            </li>
          </ul>
        </div>
      </div>
      <div class="text-sm text-gray-60 w-full text-center py-[18px]">
        {{ t('sections.footer.rights', { year: new Date().getFullYear() }) }}
      </div>
    </div>
  </footer>
</template>

<script lang="ts" setup>
import { TswIcon } from '@tsuwari/ui-components';

import NavMenu from '@/components/landing/layout/NavMenu.vue';
import type { NavMenuLocale } from '@/types/navMenu.js';
import type { SocialMediaItem } from '@/types/socialMedia.js';
import { useTranslation } from '@/utils/locales.js';

defineProps<{ menuItems: NavMenuLocale[]; socials: SocialMediaItem[] }>();

const t = useTranslation<'landing'>();
</script>
