<template>
  <header class="flex container py-3 items-center justify-between">
    <div class="inline-grid items-center grid-flow-col gap-x-[10px] p-2">
      <img src="@/assets/NewLogo.svg" />
      <span class="font-medium text-xl">Tsuwari</span>
    </div>
    <nav>
      <ul class="inline-grid grid-flow-col gap-x-2">
        <li v-for="item in menuItems" :key="item.id">
          <a :href="rt(item.href)" class="leading-tight px-3 py-[10px]">{{ rt(item.name) }}</a>
        </li>
      </ul>
    </nav>
    <div class="inline-grid grid-flow-col gap-x-3 items-center">
      <LangSelect @change="setLocale" />
      <a href="#" class="inline-flex bg-purple-60 px-4 py-[10px] rounded-md leading-tight">Login</a>
    </div>
  </header>
</template>

<script lang="ts" setup>
import { navigate } from 'vite-plugin-ssr/client/router';
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
</script>
