<script setup lang="ts">
import { NAvatar, NSelect, NSpace, NTag, NText, type SelectOption } from 'naive-ui'
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'

import type { VNodeChild } from 'vue'

import { useTwitchRewardsNew } from '@/api/twitch.ts'
import RewardFallbackImg from '@/assets/images/reward-fallback.png?url'

const props = defineProps<{
	multiple?: boolean
	clearable?: boolean
	onlyWithInput?: boolean
	placeholder?: string
}>()

const modelValue = defineModel<string | string[] | undefined | null>()

const { t } = useI18n()

const {
	data: rewardsData,
	fetching: isRewardsLoading,
	error: isRewardsError,
} = useTwitchRewardsNew()

type RewardSelectOptions = SelectOption & {
	image?: string
	color: string
}

const rewardsSelectOptions = computed(() => {
	const rewards: RewardSelectOptions[] = []
	if (!rewardsData.value?.twitchRewards) return rewards

	for (const reward of rewardsData.value.twitchRewards) {
		if (props.onlyWithInput && !reward.userInputRequired) continue

		rewards.push({
			value: reward.id,
			label: reward.title,
			image: reward.imageUrls?.at(-1),
			color: reward.backgroundColor,
			disabled: !reward.enabled,
		})
	}

	return rewards
})

function renderRewardLabel(option: RewardSelectOptions): VNodeChild {
	return h(NSpace, { align: 'center' }, {
		default: () => [
			h(NAvatar, {
				src: option.image || RewardFallbackImg,
				color: option.color,
				class: 'flex w-5 h-5 p-1',
			}),
			h(NText, {
				style: {
					color: option.disabled
						? 'var(--n-text-disabled-color)'
						: 'var(--n-text-color)',
				},
			}, { default: () => option.label }),
		],
	})
}

function renderRewardTag(props: { option: SelectOption, handleClose: () => void }): VNodeChild {
	return h(NTag, {
		bordered: false,
		closable: true,
		onClose: props.handleClose,
	},	{
		icon: () => h('img', { src: props.option.image || RewardFallbackImg, class: 'w-4 h-4 mr-1' }),
		default: () => props.option.label,
	})
}
</script>

<template>
	<NSelect
		v-model:value="modelValue"
		class="bg-background"
		:multiple="multiple"
		size="large"
		:options="rewardsSelectOptions"
		:placeholder="placeholder ?? t('events.targetTwitchReward')"
		:loading="isRewardsLoading"
		:render-label="renderRewardLabel"
		:render-tag="renderRewardTag"
		:disabled="isRewardsError !== undefined"
		:clearable="clearable"
		:virtual-scroll="false"
		filterable
	/>
</template>
