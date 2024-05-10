<script setup lang="ts">
import { useDebounce } from '@vueuse/core'
import {
	NAvatar,
	NButton,
	NForm,
	NFormItem,
	NInputNumber,
	NSelect,
	NSpace,
	NSwitch,
	NTabPane,
	NTabs,
	NText,
	type SelectOption,
	useMessage,
} from 'naive-ui'
import { computed, h, ref, toRaw, unref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import RewardsSelector from '../rewardsSelector.vue'

import type { SongRequestsSettingsOpts } from '@/gql/graphql'
import type { VNodeChild } from 'vue'

import { useCommandsApi } from '@/api/commands/commands'
import { useSongRequestsApi } from '@/api/song-requests'
import TwitchSearchUsers from '@/components/twitchUsers/multiple.vue'
import CommandList from '@/features/commands/components/list.vue'
import { SongRequestsSearchChannelOrVideoOptsType } from '@/gql/graphql'

const { t } = useI18n()

const youtubeModuleManager = useSongRequestsApi()
const youtubeModuleData = youtubeModuleManager.useSongRequestQuery()
const youtubeModuleSettings = computed(() => {
	return youtubeModuleData.data?.value?.songRequests
})
const youtubeModuleUpdater = youtubeModuleManager.useSongRequestMutation()

const formValue = ref<SongRequestsSettingsOpts>({
	enabled: false,
	acceptOnlyWhenOnline: true,
	takeSongFromDonationMessages: false,
	playerNoCookieMode: false,
	channelPointsRewardId: '',
	maxRequests: 500,
	announcePlay: true,
	neededVotesForSkip: 30,
	denyList: {
		artistsNames: [],
		songs: [],
		users: [],
		channels: [],
		words: [],
	},
	song: {
		maxLength: 10,
		minLength: 0,
		minViews: 50000,
		acceptedCategories: [],
	},
	user: {
		maxRequests: 20,
		minWatchTime: 0,
		minFollowTime: 0,
		minMessages: 0,
	},
	translations: {
		notEnabled: 'Song requests not enabled.',
		nowPlaying: 'Now playing "{{songTitle}}" {{songLink}} requested by @{{orderedByDisplayName}}',
		noText: 'You should provide text for song request.',
		acceptOnlyWhenOnline: 'Requests accepted only on online streams.',
		song: {
			notFound: 'Song not found.',
			alreadyInQueue: 'Song already in queue.',
			ageRestrictions: 'Age restriction on that song.',
			cannotGetInformation: 'Cannot get information about song.',
			live: 'Seems like that song is live, which is disallowed.',
			denied: 'That song is denied for requests.',
			requestedMessage: 'Song "{{songTitle}}" requested, queue position {{position}}. Estimated wait time before your track will be played is {{waitTime}}.',
			maximumOrdered: 'Maximum number of songs is queued ({{maximum}}).',
			minViews: 'Song {{songTitle}} ({{songViews}} views) haven\'t {{neededViews}} views for being ordered',
			maxLength: 'Maximum length of song is {{maxLength}}',
			minLength: 'Minimum length of song is {{minLength}}',
		},
		user: {
			denied: 'You are denied to request any song.',
			maxRequests: 'Maximum number of songs ordered by you ({{count}})',
			minMessages: 'You have only {{userMessages}} messages, but needed {{neededMessages}} for requesting song',
			minWatched: 'You\'ve only watched {{userWatched}} but needed {{neededWatched}} to request a song.',
			minFollow: 'You are followed for {{userFollow}} minutes, but needed {{neededFollow}} for requesting song',
		},
		channel: {
			denied: 'That channel is denied for requests.',
		},
	},
})

watch(youtubeModuleSettings, async (v) => {
	if (!v) return
	formValue.value = toRaw(v)
}, { immediate: true })

const message = useMessage()

async function save() {
	const data = unref(formValue)

	await youtubeModuleUpdater.executeMutation({ opts: data })
	message.success(t('sharedTexts.saved'))
}

function renderSelectOption(option: SelectOption & { image: string }): VNodeChild {
	return h(NSpace, {
		align: 'center',
	}, {
		default: () => [
			h(NAvatar, {
				src: option.image,
				round: true,
				class: 'flex h-5 w-5',
			}),
			h(NText, { class: 'ml-2' }, { default: () => option.label }),
		],
	})
}

const channelsSearchValue = ref('')
const channelsSearchDebounced = useDebounce(channelsSearchValue, 500)

const channelsSearchOpts = computed(() => {
	return {
		query: [...formValue!.value!.denyList!.channels, channelsSearchDebounced.value],
		type: SongRequestsSearchChannelOrVideoOptsType.Channel,
	}
})
const channelsSearch = youtubeModuleManager.useYoutubeVideoOrChannelSearch(channelsSearchOpts)

const channelsOptions = computed(() => {
	return channelsSearch?.data.value?.songRequestsSearchChannelOrVideo.items.map((channel) => {
		return {
			label: channel.title,
			value: channel.id,
			image: channel.thumbnail,
		}
	})
})

const songsSearchValue = ref('')
const songsSearchDebounced = useDebounce(songsSearchValue, 500)

const songsSearchOpts = computed(() => {
	return {
		query: [...formValue!.value!.denyList!.songs, songsSearchDebounced.value],
		type: SongRequestsSearchChannelOrVideoOptsType.Video,
	}
})

const songsSearch = youtubeModuleManager.useYoutubeVideoOrChannelSearch(songsSearchOpts)
const songsSearchOptions = computed(() => {
	return songsSearch?.data.value?.songRequestsSearchChannelOrVideo.items.map((channel) => {
		return {
			label: channel.title,
			value: channel.id,
			image: channel.thumbnail,
		}
	}) ?? []
})

const commandsManager = useCommandsApi()
const { data: commands } = commandsManager.useQueryCommands()
const srCommands = computed(() => {
	return commands.value?.commands.filter((c) => c.module === 'SONGS' && c.defaultName !== 'song') ?? []
})
</script>

<template>
	<NForm>
		<NTabs
			type="line"
			animated
		>
			<NTabPane name="general" :tab="t('songRequests.tabs.general')">
				<NSpace vertical>
					<NSpace justify="space-between">
						<NText>{{ t('sharedTexts.enabled') }}</NText>
						<NSwitch v-model:value="formValue.enabled" />
					</NSpace>

					<NSpace justify="space-between">
						<NText>{{ t('songRequests.settings.takeSongFromDonationMessage') }}</NText>
						<NSwitch v-model:value="formValue.takeSongFromDonationMessages" />
					</NSpace>

					<NSpace justify="space-between">
						<NText>{{ t('songRequests.settings.onlineOnly') }}</NText>
						<NSwitch v-model:value="formValue.acceptOnlyWhenOnline" />
					</NSpace>

					<NSpace justify="space-between">
						<NText>{{ t('songRequests.settings.announcePlay') }}</NText>
						<NSwitch v-model:value="formValue.announcePlay" />
					</NSpace>

					<div>
						<NSpace justify="space-between">
							<NText>{{ t('songRequests.settings.playerNoCookieMode') }}</NText>
							<NSwitch v-model:value="formValue.playerNoCookieMode" />
						</NSpace>
						<NText class="text-xs mt-1">
							{{ t('songRequests.settings.playerNoCookieModeDescription') }}
						</NText>
					</div>

					<NFormItem
						:label="t('songRequests.settings.neededPercentageForskip')"
						path="neededVotesVorSkip"
					>
						<NInputNumber v-model:value="formValue.neededVotesForSkip" />
					</NFormItem>

					<NFormItem
						:label="t('songRequests.settings.channelReward')"
						path="channelPointsRewardId"
					>
						<RewardsSelector
							v-model="formValue.channelPointsRewardId"
							only-with-input
							clearable
						/>
					</NFormItem>

					<NFormItem
						:label="t('songRequests.settings.deniedChannels')"
						path="channelPointsRewardId"
					>
						<NSelect
							v-model:value="formValue.denyList!.channels"
							:loading="channelsSearch.fetching.value"
							remote
							filterable
							:options="channelsOptions"
							clearable
							multiple
							:clear-filter-after-select="false"
							:render-label="renderSelectOption"
							@search="(v) => channelsSearchValue = v"
						/>
					</NFormItem>

					<NFormItem
						:label="t('songRequests.settings.deniedWords')"
						path="channelPointsRewardId"
					>
						<NSelect
							v-model:value="formValue.denyList!.words"
							filterable
							multiple
							clearable
							tag
						/>
					</NFormItem>
				</NSpace>
			</NTabPane>

			<NTabPane name="commands" :tab="t('commands.name')">
				<CommandList class="mb-2" :commands="srCommands" />
			</NTabPane>

			<NTabPane name="users" :tab="t('songRequests.tabs.users')">
				<NFormItem :label="t('songRequests.settings.users.maxRequests')" path="user.maxRequests">
					<NInputNumber v-model:value="formValue.user!.maxRequests" :min="0" :max="1000" />
				</NFormItem>
				<NFormItem
					:label="t('songRequests.settings.users.minimalWatchTime')"
					path="user.minWatchTime"
				>
					<NInputNumber v-model:value="formValue.user!.minWatchTime" :min="0" :max="999999999" />
				</NFormItem>
				<NFormItem
					:label="t('songRequests.settings.users.minimalMessages')"
					path="user.minMessages"
				>
					<NInputNumber v-model:value="formValue.user!.minMessages" :min="0" :max="999999999" />
				</NFormItem>
				<NFormItem
					:label="t('songRequests.settings.users.minimalFollowTime')"
					path="user.minFollowTime"
				>
					<NInputNumber
						v-model:value="formValue.user!.minFollowTime" :min="0"
						:max="99999999999999"
					/>
				</NFormItem>

				<NFormItem :label="t('songRequests.settings.deniedUsers')">
					<TwitchSearchUsers v-model="formValue.denyList!.users" />
				</NFormItem>
			</NTabPane>

			<NTabPane name="songs" :tab="t('songRequests.tabs.songs')">
				<NFormItem :label="t('songRequests.settings.songs.maxRequests')">
					<NInputNumber v-model:value="formValue.maxRequests" :min="0" :max="99999999999999" />
				</NFormItem>
				<NFormItem :label="t('songRequests.settings.songs.minLength')">
					<NInputNumber v-model:value="formValue.song!.minLength" :min="0" :max="999999" />
				</NFormItem>
				<NFormItem :label="t('songRequests.settings.songs.maxLength')">
					<NInputNumber v-model:value="formValue.song!.maxLength" :min="0" :max="999999" />
				</NFormItem>
				<NFormItem :label="t('songRequests.settings.songs.minViews')">
					<NInputNumber v-model:value="formValue.song!.minViews" :min="0" :max="99999999999999" />
				</NFormItem>
				<NFormItem :label="t('songRequests.settings.deniedSongs')">
					<NSelect
						v-model:value="formValue.denyList!.songs"
						:loading="songsSearch.fetching.value"
						remote
						filterable
						multiple
						:options="songsSearchOptions"
						:render-label="renderSelectOption"
						clearable
						@search="(v) => songsSearchValue = v"
					/>
				</NFormItem>
			</NTabPane>
		</NTabs>

		<NButton secondary block type="success" @click="save">
			{{ t('sharedButtons.save') }}
		</NButton>
	</NForm>
</template>
