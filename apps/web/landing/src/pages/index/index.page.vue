<template>
  <Header :menuItems="data.navMenuItems" />
  <main>
    {{ t('hello') }}
    <div class="locale-changer">
      <select @change="(e) => setLocale(e.target!.value)">
        <option v-for="l in locales" :key="`locale-${l}`" :value="l">
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
  </main>
  <Footer :menuItems="data.navMenuItems" :socials="data.socials" />
</template>

<script lang="ts" setup>
import { nextTick } from 'vue';
import { useI18n } from 'vue-i18n';

import FeaturesBlock from '@/components/sections/FeaturesBlock.vue';
import FirstScreen from '@/components/sections/FirstScreen.vue';
import Footer from '@/components/sections/Footer.vue';
import Header from '@/components/sections/Header.vue';
import IntegrationsBlock from '@/components/sections/IntegrationsBlock.vue';
import PricingBlock from '@/components/sections/PricingBlock.vue';
import ReviewsBlock from '@/components/sections/ReviewsBlock.vue';
import StatsLine from '@/components/sections/StatsLine.vue';
import TeamBlock from '@/components/sections/TeamBlock.vue';
import * as data from '@/data';
import type { Locale } from '@/types/locale.js';
import { locales } from '@/utils/locales.js';

const { setLocaleMessage, t, locale } = useI18n();

async function setLocale(l: Locale) {
  const messages = await import(`../../locales/landing/${l}.json`);

  setLocaleMessage(l, messages.default);
  locale.value = l;

  return nextTick();
}


</script>
