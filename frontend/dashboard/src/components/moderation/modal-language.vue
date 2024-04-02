<script setup lang="ts">
import { NTransfer, NDivider } from 'naive-ui';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useEditableItem } from './helpers.js';

import { useModerationAvailableLanguages } from '@/api';

const { editableItem } = useEditableItem();

const { data: availableLanguages } = useModerationAvailableLanguages();

const transferOptions = computed(() => {
	return availableLanguages?.value?.langs.map(l => ({
		label: l.name,
		value: l.code.toString(),
	})) ?? [];
});

const { t } = useI18n();
</script>

<template>
	<div>
		<n-transfer
			ref="transfer"
			v-model:value="editableItem!.data!.deniedChatLanguages"
			:show-selected="false"
			virtual-scroll
			:options="transferOptions"
			source-filterable
			target-filterable
			:source-title="t('moderation.types.language.allowedLanguages')"
			:target-title="t('moderation.types.language.disallowedLanguages')"
			source-filter-placeholder="Search"
			target-filter-placeholder="Search"
		/>

		<n-divider class="m-0 p-0" />
	</div>
</template>
