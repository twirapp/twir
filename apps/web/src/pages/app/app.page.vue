<template>
  <span>{{ t('app.tabs.commands') }}</span>
  <div class="locale-changer">
    <select
      @change="(e) => setLocale(
        (e.target! as HTMLSelectElement).value as Locale)
      "
    >
      <option v-for="l in locales" :key="`locale-${l}`" :value="l" :selected="l == locale">
        {{ l }}
      </option>
    </select>
  </div>
  <router-view v-slot="{ Component }">
    <Suspense>
      <component :is="Component" />
    </Suspense>
  </router-view>
  <router-link to="/dashboard">
    Dashboard
  </router-link>
  <router-link to="/commands">
    Commands
  </router-link>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n';

import type { Locale } from '../../types/locale.js';
import { locales } from '../../utils/locales.js';

const { t, setLocaleMessage, locale } = useI18n();

async function setLocale(l: Locale) {
  const message = (await import(`../../locales/app/${l}.json`)).default;

  setLocaleMessage('ru', message);

  locale.value = l;
}
</script>
