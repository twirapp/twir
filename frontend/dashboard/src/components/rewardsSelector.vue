<script setup lang="ts">
import { SelectOption, NSpace, NAvatar, NText, NSelect } from 'naive-ui';
import { computed, VNodeChild, h } from 'vue';
import { useI18n } from 'vue-i18n';

import { useTwitchRewards } from '@/api';

const props = defineProps<{
	multiple?: boolean
	clearable?: boolean
	onlyWithInput?: boolean
}>();

// eslint-disable-next-line no-undef
const modelValue = defineModel<string | string[]>();

const { t } = useI18n();

const {
	data: rewardsData,
	isLoading: isRewardsLoading,
	isError: isRewardsError,
} = useTwitchRewards();

const rewardsSelectOptions = computed(() => {
	return rewardsData.value?.rewards.map(r => ({
		value: r.id,
		label: r.title,
		image: r.image?.url4X,
		disabled: props.onlyWithInput ?? false,
	})) ?? [];
});

const renderRewardTag = (option: SelectOption & { image?: string }): VNodeChild => {
	return h(NSpace, { align: 'center' }, {
		default: () => [
			h(NAvatar, { src: option.image, round: true, size: 'small', style: 'display: flex;' }),
			h(NText, {}, { default: () => option.label }),
		],
	});
};
</script>

<template>
	<n-select
		v-model:value="modelValue"
		:multiple="multiple"
		size="large"
		:options="rewardsSelectOptions"
		:placeholder="t('events.targetTwitchReward')"
		:loading="isRewardsLoading"
		:render-label="renderRewardTag"
		:disabled="isRewardsError"
		:clearable="clearable"
		filterable
	/>
</template>
