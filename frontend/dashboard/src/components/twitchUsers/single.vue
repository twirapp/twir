<script setup lang='ts'>
import { refDebounced } from '@vueuse/core';
import { NSelect, NTag, NAvatar } from 'naive-ui';
import { computed, ref, watch, h } from 'vue';

import { useTwitchSearchChannels, useTwitchGetUsers } from '@/api/index.js';

// eslint-disable-next-line no-undef
const userId = defineModel<string>({ default: '' });

const props = defineProps<{
	initialUserId?: string;
}>();

const getUsers = useTwitchGetUsers({
	ids: [props.initialUserId ?? ''],
});

const userName = ref<string>('');
const userNameDebounced = refDebounced(userName, 1000);
const twitchSearch = useTwitchSearchChannels(userNameDebounced);
watch(userNameDebounced, () => {
	twitchSearch.refetch();
});

const options = computed(() => {
	const searchUsers = twitchSearch.data.value?.channels ?? [];
	const initialUsers = getUsers.data.value?.users ?? [];

	return [
		...searchUsers.map((channel) => ({
			label: channel.login === channel.displayName.toLowerCase()
				? channel.displayName
				: `${channel.login} (${channel.displayName})`,
			value: channel.id,
			profileImageUrl: channel.profileImageUrl,
		})).filter((channel) => !initialUsers.find((u) => u.id === channel.value)),
		...initialUsers.map((user) => ({
			label: user.login === user.displayName.toLowerCase()
				? user.displayName
				: `${user.login} (${user.displayName})`,
			value: user.id,
			profileImageUrl: user.profileImageUrl,
		})),
	];
});

function handleSearch(query: string) {
	userName.value = query;
}

type Option = {
	label: string;
	value: string;
	profileImageUrl?: string;
};

const renderMultipleSelectTag = ({ option }: {
	option: Option;
	handleClose: () => void;
}) => {
	return !option.label ? h('p') : h(
		NTag,
		{
			style: {
				padding: '0 6px 0 4px',
			},
			round: true,
		},
		{
			default: () =>
				h(
					'div',
					{
						style: {
							display: 'flex',
							alignItems: 'center',
						},
					},
					[
						h(NAvatar, {
							src: option.profileImageUrl,
							round: true,
							size: 22,
							style: {
								marginRight: '4px',
							},
						}),
						option.label as string,
					],
				),
		},
	);
};

const renderLabel = (option: Option | undefined) => {
	return !option ? null : h(
		'div',
		{
			style: {
				display: 'flex',
				alignItems: 'center',
			},
		},
		[
			h(NAvatar, {
				src: option.profileImageUrl,
				round: true,
				size: 'small',
			}),
			h(
				'div',
				{
					style: {
						marginLeft: '12px',
						padding: '4px 0',
					},
				},
				[
					h('div', null, option.label),
				],
			),
		],
	);
};
</script>

<template>
	<n-select
		v-model:value="userId"
		filterable
		placeholder="Search users..."
		:options="options"
		:loading="twitchSearch.isLoading.value"
		clearable
		remote
		:clear-filter-after-select="true"
		:render-label="renderLabel as any"
		:render-tag="renderMultipleSelectTag as any"
		@search="handleSearch"
	/>
</template>

<style scoped lang='postcss'>

</style>
