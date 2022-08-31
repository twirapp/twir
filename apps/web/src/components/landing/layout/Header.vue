<template>
  <header
    :class="`
      header
      ${isHeaderScrolled ? 'header-bb' : ''}
    `"
  >
    <div class="flex container py-3 items-center justify-between">
      <div class="flex-1 flex">
        <div class="mr-auto">
          <a class="inline-grid items-center grid-flow-col gap-x-[10px] p-2" href="#">
            <img src="@/assets/NewLogo.svg" />
            <span class="font-medium text-xl">Tsuwari</span>
          </a>
        </div>
      </div>
      <nav>
        <ul class="inline-grid grid-flow-col gap-x-2">
          <li v-for="item in menuItems" :key="item.id">
            <a
              :href="rt(item.href)"
              class="
                leading-tight
                px-3
                py-[10px]
                text-gray-70
                transition-colors
                hover:text-white-100
              "
            >
              {{ rt(item.name) }}
            </a>
          </li>
        </ul>
      </nav>
      <div class="flex-1 flex">
        <div class="inline-grid grid-flow-col gap-x-3 items-center ml-auto">
          <LangSelect @change="setLocale" />
          <a
            href="#"
            class="
              inline-flex
              bg-purple-60
              px-[13px]
              py-[7px]
              rounded-md
              leading-tight
              hover:bg-opacity-20
              hover:border-opacity-50
              hover:text-purple-95
              border-2 border-opacity-0 border-purple-70
              transition-colors
            "
          >
            Login
          </a>
        </div>
      </div>
    </div>
  </header>
</template>

<script lang="ts" setup>
import { useWindowScroll } from '@vueuse/core';
import { navigate } from 'vite-plugin-ssr/client/router';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import LangSelect from '@/components/LangSelect/LangSelect.vue';
import type { Locale } from '@/types/locale.js';
import type { NavMenuItem } from '@/types/navMenu';
import { loadLocaleMessages } from '@/utils/locales.js';

defineProps<{ menuItems: NavMenuItem[] }>();

const { setLocaleMessage, locale: i18nLocale, rt } = useI18n();

async function setLocale(locale: Locale) {
  const messages = await loadLocaleMessages('landing', locale);

  setLocaleMessage<any>(locale, messages);
  i18nLocale.value = locale;

  navigate(`/${locale}`);
}

const { y } = useWindowScroll();

const isHeaderScrolled = computed(() => {
  return y.value > 70;
});
</script>

<style lang="postcss" scoped>
.header {
  @apply sticky
    w-full
    left-0
    right-0
    top-0
    mx-auto
    z-20
    bg-black-10 bg-opacity-0
    border-b border-opacity-0 border-black-20
    backdrop-blur-sm backdrop-saturate-[180%];

  transition: border-bottom-color 0.3s ease, background 0.3s ease;
}

.header-bb {
  @apply border-opacity-80 bg-opacity-80;
}
</style>
