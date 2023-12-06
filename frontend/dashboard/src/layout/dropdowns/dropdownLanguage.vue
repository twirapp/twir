<script setup lang="ts">
import { useLocalStorage } from '@vueuse/core';
import { NButton, NDropdown } from 'naive-ui';
import { defineAsyncComponent } from 'vue';
import { useI18n } from 'vue-i18n';

const { t, locale, availableLocales } = useI18n();

const localStorageLocale = useLocalStorage('twirLocale', 'en');

const renderFlagIcon = (code: string) => defineAsyncComponent(() => import(`../../assets/icons/flags/${code}.svg?component`));
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
			localStorageLocale = l
		}"
	>
		<n-button circle quaternary style="padding: 5px; font-size: 25px">
			<component :is="renderFlagIcon(localStorageLocale)" style="width: 35px; height: 50px;" />
		</n-button>
	</n-dropdown>
</template>
