<script setup lang="ts">
import { IconWorld } from '@tabler/icons-vue';
import { useLocalStorage } from '@vueuse/core';
import { NButton, NDropdown } from 'naive-ui';
import { useI18n } from 'vue-i18n';

import { AVAILABLE_LOCALES } from '@/plugins/i18n';

const { locale } = useI18n();

const currentLocale = useLocalStorage<string>('twirLocale', 'en');
</script>

<template>
	<n-dropdown
		trigger="click"
		:options="AVAILABLE_LOCALES.map((locale) => ({
			title: locale.name,
			key: locale.code,
		}))"
		size="medium"
		@select="(newLocale) => {
			locale = newLocale
			currentLocale = newLocale
		}"
	>
		<n-button quaternary class="!text-[16px] !p-[5px]">
			<div class="flex gap-2 items-center">
				<IconWorld />
				{{ currentLocale.toUpperCase() }}
			</div>
		</n-button>
	</n-dropdown>
</template>
