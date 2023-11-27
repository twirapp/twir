<script setup lang="ts">
import type { ItemWithId } from '@twir/grpc/generated/api/api/moderation';
import { NTransfer, NDivider } from 'naive-ui';
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useModerationAvailableLanguages } from '@/api';

defineProps<{
	item: ItemWithId
}>();

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
			v-model:value="item.data!.deniedChatLanguages"
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

		<n-divider style="margin: 0; padding: 0" />
	</div>
</template>
