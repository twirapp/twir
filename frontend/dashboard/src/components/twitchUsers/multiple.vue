<script setup lang="ts">
import type { Channel, TwitchSearchChannelsRequest, TwitchUser } from '@twir/api/messages/twitch/twitch';
import { refDebounced } from '@vueuse/core';
import { NSelect, NTag, NAvatar } from 'naive-ui';
import { computed, ref, h } from 'vue';

import { useTwitchSearchChannels, useTwitchGetUsers } from '@/api/index.js';
import { resolveUserName } from '@/helpers';

// eslint-disable-next-line no-undef
const usersIds = defineModel<string[]>({ default: [] });
defineProps<{
	max?: number
}>();

const getUsers = useTwitchGetUsers({
	ids: usersIds,
});

const userName = ref('');
const userNameDebounced = refDebounced(userName, 500);

const searchParams = computed<TwitchSearchChannelsRequest>(() => ({
	query: userNameDebounced.value,
	twirOnly: false,
}));
const twitchSearch = useTwitchSearchChannels(searchParams);

function mapOptions(users: (TwitchUser | Channel)[] ) {
	return users.map((user) => ({
		label: resolveUserName(user.login, user.displayName),
		value: user.id,
		profileImageUrl: user.profileImageUrl,
	}));
}

const options = computed(() => {
	const searchUsers = twitchSearch.data.value?.channels ?? [];
	const initialUsers = getUsers.data.value?.users ?? [];

	return [
		...mapOptions(searchUsers)
			.filter((channel) => !initialUsers.find((user) => user.id === channel.value)),
		...mapOptions(initialUsers),
	];
});


function handleSearch(query: string) {
	userName.value = query;
}

type Option = {
	label: string;
	value: string;
	profileImageUrl: string;
};

const renderMultipleSelectTag = ({ option, handleClose }: {
	option: Option;
	handleClose: () => void;
}) => {
	return h(
		NTag,
		{
			class: 'pr-1.5 pl-1',
			round: true,
			closable: true,
			onClose: (e) => {
				e.stopPropagation();
				handleClose();
			},
		},
		{
			default: () =>
				h(
					'div',
					{ class: 'flex items-center' },
					[
						h(NAvatar, {
							src: option.profileImageUrl,
							round: true,
							size: 22,
							class: 'mr-1',
						}),
						option.label as string,
					],
				),
		},
	);
};

const renderLabel = (option: Option) => {
	return h(
		'div',
		{ class: 'flex items-center' },
		[
			h(NAvatar, {
				src: option.profileImageUrl,
				round: true,
				size: 'small',
			}),
			h(
				'div',
				{ class: 'ml-3 py-1' },
				[h('div', null, option.label)],
			),
		],
	);
};
</script>

<template>
	<n-select
		v-model:value="usersIds"
		multiple
		:filterable="max ? usersIds.length !== max : true"
		placeholder="Search users..."
		:options="options"
		:loading="twitchSearch.isLoading.value"
		clearable
		remote
		:clear-filter-after-select="true"
		:render-label="renderLabel as any"
		:render-tag="renderMultipleSelectTag as any"
		@search="handleSearch"
		@update-value="userName = ''"
	/>
</template>
