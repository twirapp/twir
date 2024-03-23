<script setup lang="ts">
import { IconWorld } from '@tabler/icons-vue';
import { useLocalStorage } from '@vueuse/core';
import { NButton, NDropdown } from 'naive-ui';
import { useI18n } from 'vue-i18n';

const { t, locale, availableLocales } = useI18n();

const currentLocale = useLocalStorage<string>('twirLocale', 'en');
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
		<n-button quaternary style="padding: 5px; font-size: 16px">
			<div class="flex gap-2 items-center">
				<IconWorld />
				{{ currentLocale.toUpperCase() }}
			</div>
		</n-button>
	</n-dropdown>
</template>
