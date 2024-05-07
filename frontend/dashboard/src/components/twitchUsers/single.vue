<script setup lang='ts'>
import { refDebounced } from '@vueuse/core'
import { NAvatar, NSelect, NTag } from 'naive-ui'
import { computed, h, ref } from 'vue'

import type { Channel, TwitchSearchChannelsRequest, TwitchUser } from '@twir/api/messages/twitch/twitch'

import { useTwitchGetUsers, useTwitchSearchChannels } from '@/api/index.js'
import { resolveUserName } from '@/helpers'

const props = withDefaults(defineProps<{
	initialUserId?: string | null
	twirOnly?: boolean
}>(), {
	twirOnly: false,
	initialUserId: null,
})

const userId = defineModel<string | null>({ default: null })

const getUsers = useTwitchGetUsers({
	ids: [props.initialUserId ?? ''],
})

const userName = ref('')
const userNameDebounced = refDebounced(userName, 500)

const searchParams = computed<TwitchSearchChannelsRequest>(() => ({
	query: userNameDebounced.value,
	twirOnly: props.twirOnly,
}))
const twitchSearch = useTwitchSearchChannels(searchParams)

function mapOptions(users: (TwitchUser | Channel)[]) {
	return users.map((user) => ({
		label: resolveUserName(user.login, user.displayName),
		value: user.id,
		profileImageUrl: user.profileImageUrl,
	}))
}

const options = computed(() => {
	const searchUsers = twitchSearch.data.value?.channels ?? []
	const initialUsers = getUsers.data.value?.users ?? []

	return [
		...mapOptions(searchUsers)
			.filter((channel) => !initialUsers.find((user) => user.id === channel.value)),
		...mapOptions(initialUsers),
	]
})

function handleSearch(query: string) {
	userName.value = query
}

interface Option {
	label: string
	value: string
	profileImageUrl?: string
}

function renderMultipleSelectTag({ option }: {
	option: Option
	handleClose: () => void
}) {
	return !option.label
		? h('p')
		: h(
			NTag,
			{
				class: 'pr-1.5 pl-1',
				round: true,
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
		)
}

function renderLabel(option: Option | undefined) {
	return !option
		? null
		: h(
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
		)
}
</script>

<template>
	<NSelect
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
