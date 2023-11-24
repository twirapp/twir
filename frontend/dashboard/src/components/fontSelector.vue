<script setup lang="ts">
import { NTreeSelect } from 'naive-ui';
import { useThemeVars, type TreeSelectOption } from 'naive-ui';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useGoogleFontsList } from '@/api';

withDefaults(defineProps<{
	clearable: boolean,
}>(), {
	clearable: false,
});

// eslint-disable-next-line no-undef
const modelValue = defineModel<string | string[]>();

const { t } = useI18n();
const themeVars = useThemeVars();

const fontSelectOptions = computed<TreeSelectOption[]>(() => {
  return googleFonts?.value?.fonts
      .map((f) => {
        const option: TreeSelectOption = {
          label: f.family,
          children: f.files.map((c) => ({
            label: `${f.family}:${c.name}`,
            key: `${f.family}:${c.name}`,
          })),
          key: f.family,
        };

        return option;
      }) ?? [];
});

const {
  data: googleFonts,
  isError: isGoogleFontsError,
  isLoading: isGoogleFontsLoading,
} = useGoogleFontsList();
</script>

<template>
	<n-tree-select
		v-model:value="modelValue"
		filterable
		:options="fontSelectOptions"
		:loading="isGoogleFontsLoading"
		:disabled="isGoogleFontsError"
		check-strategy="child"
		:clearable="clearable"
	>
		<template #action>
			{{ t('components.fontFamily.description') }}
			<a
				class="action-link"
				href="https://fonts.google.com/"
				target="_blank"
				:style="{ color: themeVars.successColor }"
			>
				Preview Google Fonts
			</a>
		</template>
	</n-tree-select>
</template>
