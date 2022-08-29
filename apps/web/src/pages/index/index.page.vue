<template>
  {{ t('landing.hello') }}
  <div class="locale-changer">
    <select @change="(e) => setLocale((e.target! as HTMLSelectElement).value as Locale)">
      <option v-for="l in locales" :key="`locale-${l}`" :value="l" :selected="l == locale">
        {{ l }}
      </option>
    </select>
  </div>
  <FirstScreen />
  <StatsLine :stats="data.stats" />
  <FeaturesBlock :features="data.features" />
  <IntegrationsBlock />
  <ReviewsBlock :reviews="data.reviews" />
  <TeamBlock :teamMembers="data.teamMembers" />
  <PricingBlock :planColorThemes="data.planColorThemes" :pricePlans="data.pricePlans" />
</template>

<script lang="ts" setup>
import { navigate } from 'vite-plugin-ssr/client/router';
import { useI18n } from 'vue-i18n';

import FeaturesBlock from '@/components/sections/FeaturesBlock.vue';
import FirstScreen from '@/components/sections/FirstScreen.vue';
import IntegrationsBlock from '@/components/sections/IntegrationsBlock.vue';
import PricingBlock from '@/components/sections/PricingBlock.vue';
import ReviewsBlock from '@/components/sections/ReviewsBlock.vue';
import StatsLine from '@/components/sections/StatsLine.vue';
import TeamBlock from '@/components/sections/TeamBlock.vue';
import * as data from '@/data';
import { usePageContext } from '@/hooks/usePageContext.js';
import type { Locale } from '@/types/locale.js';
import { loadLocaleMessages, locales } from '@/utils/locales.js';

const { t, setLocaleMessage, locale: i18nLocale } = useI18n();

const { locale } = usePageContext();

async function setLocale(l: Locale) {
  const messages = await loadLocaleMessages('landing', l);
  
  setLocaleMessage<any>(l, messages);
  i18nLocale.value = l;

  navigate(`/${l}`);
}

</script>
