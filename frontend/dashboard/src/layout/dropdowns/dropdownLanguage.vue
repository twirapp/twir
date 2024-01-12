<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core';
import { NButton, NDropdown } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import LanguageFlag, { type Locale } from './languageFlag.vue';

const { t, locale, availableLocales } = useI18n();

const currentLocale = useLocalStorage<Locale>('twirLocale', 'en');
</script>

<template>
	<n-dropdown
		trigger="click"
		:options="availableLocales.map(l => ({
			title: t('languageName', {}, { locale: l }),
			key: l as string,
		}))"
		size="medium"
		@select="(l) => {
			locale = l
			currentLocale = l
		}"
	>
		<n-button circle quaternary style="padding: 5px; font-size: 25px">
			<language-flag :locale="currentLocale" />
		</n-button>
	</n-dropdown>
</template>
