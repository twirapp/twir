<script setup lang="ts">
import type { YouTubeSettings } from '@twir/api/messages/modules_sr/modules_sr';
import { useDebounce } from '@vueuse/core';
import {
	type SelectOption,
	NTabs,
	NTabPane,
	NSpace,
	NSwitch,
	NText,
	NInputNumber,
	NForm,
	NFormItem,
	NSelect,
	NAvatar,
	NButton,
	useMessage,
} from 'naive-ui';
import { ref, computed, VNodeChild, h, watch, unref, toRaw } from 'vue';
import { useI18n } from 'vue-i18n';


import RewardsSelector from '../rewardsSelector.vue';

import {
	useCommandsManager,
	useYoutubeVideoOrChannelSearch,
	YoutubeSearchType,
} from '@/api/index.js';
import { useYoutubeModuleSettings } from '@/api/index.js';
import CommandList from '@/components/commands/list.vue';
import TwitchSearchUsers from '@/components/twitchUsers/multiple.vue';

const { t } = useI18n();

const youtubeModuleManager = useYoutubeModuleSettings();
const youtubeModuleData = youtubeModuleManager.getAll();
const youtubeModuleSettings = computed(() => {
	return youtubeModuleData.data?.value?.data ?? null;
});
const youtubeModuleUpdater = youtubeModuleManager.update();

const formValue = ref<YouTubeSettings>({
	enabled: false,
	acceptOnlyWhenOnline: true,
	takeSongFromDonationMessages: false,
	playerNoCookieMode: false,
	channelPointsRewardId: '',
	maxRequests: 500,
	announcePlay: true,
	neededVotesVorSkip: 30,
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
});

watch(youtubeModuleSettings, async (v) => {
	if (!v) return;
	formValue.value = toRaw(v);
}, { immediate: true });


const message = useMessage();

async function save() {
	const data = unref(formValue);

	await youtubeModuleUpdater.mutateAsync({ data });
	message.success(t('sharedTexts.saved'));
}

const renderSelectOption = (option: SelectOption & { image: string }): VNodeChild => {
	return h(NSpace,
		{
			align: 'center',
		},
		{
			default: () => [
				h(NAvatar, {
					src: option.image,
					round: true,
					style: 'height: 25px; width: 25px; display: flex',
				}),
				h(NText, {
					style: {
						marginLeft: '0.5rem',
					},
				}, {
					default: () => option.label,
				}),
			],
		});
};

const channelsSearchValue = ref('');
const channelsSearchDebounced = useDebounce(channelsSearchValue, 500);

const channelsIds = computed(() => {
	return [...formValue!.value!.denyList!.channels, channelsSearchDebounced.value];
});
const selectedChannels = useYoutubeVideoOrChannelSearch(
	channelsIds,
	YoutubeSearchType.Channel,
);

const channelsOptions = computed(() => {
	return selectedChannels?.data?.value?.items.map((channel) => {
		return {
			label: channel.title,
			value: channel.id,
			image: channel.thumbnail,
		};
	});
});

const songsSearchValue = ref('');
const songsSearchDebounced = useDebounce(songsSearchValue, 500);

const songsIds = computed(() => {
	return [...formValue!.value!.denyList!.songs, songsSearchDebounced.value];
});
const songsSearch = useYoutubeVideoOrChannelSearch(
	songsIds,
	YoutubeSearchType.Video,
);
const songsSearchOptions = computed(() => {
	return songsSearch.data?.value?.items.map((channel) => {
		return {
			label: channel.title,
			value: channel.id,
			image: channel.thumbnail,
		};
	}) ?? [];
});

const { data: allCommands } = useCommandsManager().getAll({});
const srCommands = computed(() => {
	return allCommands.value?.commands.filter((c) => c.module === 'SONGS' && c.defaultName !== 'song') ?? [];
});
</script>

<template>
	<n-form>
		<n-tabs
			type="line"
			animated
		>
			<n-tab-pane name="general" :tab="t('songRequests.tabs.general')">
				<n-space vertical>
					<n-space justify="space-between">
						<n-text>{{ t('sharedTexts.enabled') }}</n-text>
						<n-switch v-model:value="formValue.enabled" />
					</n-space>

					<n-space justify="space-between">
						<n-text>{{ t('songRequests.settings.takeSongFromDonationMessage') }}</n-text>
						<n-switch v-model:value="formValue.takeSongFromDonationMessages" />
					</n-space>

					<n-space justify="space-between">
						<n-text>{{ t('songRequests.settings.onlineOnly') }}</n-text>
						<n-switch v-model:value="formValue.acceptOnlyWhenOnline" />
					</n-space>

					<n-space justify="space-between">
						<n-text>{{ t('songRequests.settings.announcePlay') }}</n-text>
						<n-switch v-model:value="formValue.announcePlay" />
					</n-space>

					<n-space justify="space-between" style="margin-bottom: 12px;">
						<n-text>{{ t('songRequests.settings.playerNoCookieMode') }}</n-text>
						<n-switch v-model:value="formValue.playerNoCookieMode" />
						<n-text style="font-size: 12px; margin-top: 4px;">
							{{ t('songRequests.settings.playerNoCookieModeDescription') }}
						</n-text>
					</n-space>

					<n-form-item
						:label="t('songRequests.settings.neededPercentageForskip')"
						path="neededVotesVorSkip"
					>
						<n-input-number v-model:value="formValue.neededVotesVorSkip" />
					</n-form-item>

					<n-form-item
						:label="t('songRequests.settings.channelReward')"
						path="channelPointsRewardId"
					>
						<rewards-selector
							v-model="formValue.channelPointsRewardId"
							only-with-input
							clearable
						/>
					</n-form-item>

					<n-form-item
						:label="t('songRequests.settings.deniedChannels')"
						path="channelPointsRewardId"
					>
						<n-select
							v-model:value="formValue.denyList!.channels"
							:loading="selectedChannels.isLoading.value"
							remote
							filterable
							:options="channelsOptions"
							clearable
							multiple
							:clear-filter-after-select="false"
							:render-label="renderSelectOption"
							@search="(v) => channelsSearchValue = v"
						/>
					</n-form-item>

					<n-form-item
						:label="t('songRequests.settings.deniedWords')"
						path="channelPointsRewardId"
					>
						<n-select
							v-model:value="formValue.denyList!.words"
							filterable
							multiple
							clearable
							tag
						/>
					</n-form-item>
				</n-space>
			</n-tab-pane>

			<n-tab-pane name="commands" :tab="t('commands.name')">
				<CommandList class="mb-2" :commands="srCommands" />
			</n-tab-pane>

			<n-tab-pane name="users" :tab="t('songRequests.tabs.users')">
				<n-form-item :label="t('songRequests.settings.users.maxRequests')" path="user.maxRequests">
					<n-input-number v-model:value="formValue.user!.maxRequests" :min="0" :max="1000" />
				</n-form-item>
				<n-form-item
					:label="t('songRequests.settings.users.minimalWatchTime')"
					path="user.minWatchTime"
				>
					<n-input-number v-model:value="formValue.user!.minWatchTime" :min="0" :max="999999999" />
				</n-form-item>
				<n-form-item
					:label="t('songRequests.settings.users.minimalMessages')"
					path="user.minMessages"
				>
					<n-input-number v-model:value="formValue.user!.minMessages" :min="0" :max="999999999" />
				</n-form-item>
				<n-form-item
					:label="t('songRequests.settings.users.minimalFollowTime')"
					path="user.minFollowTime"
				>
					<n-input-number
						v-model:value="formValue.user!.minFollowTime" :min="0"
						:max="99999999999999"
					/>
				</n-form-item>

				<n-form-item :label="t('songRequests.settings.deniedUsers')">
					<twitch-search-users v-model="formValue.denyList!.users" />
				</n-form-item>
			</n-tab-pane>

			<n-tab-pane name="songs" :tab="t('songRequests.tabs.songs')">
				<n-form-item :label="t('songRequests.settings.songs.maxRequests')">
					<n-input-number v-model:value="formValue.maxRequests" :min="0" :max="99999999999999" />
				</n-form-item>
				<n-form-item :label="t('songRequests.settings.songs.minLength')">
					<n-input-number v-model:value="formValue.song!.minLength" :min="0" :max="999999" />
				</n-form-item>
				<n-form-item :label="t('songRequests.settings.songs.maxLength')">
					<n-input-number v-model:value="formValue.song!.maxLength" :min="0" :max="999999" />
				</n-form-item>
				<n-form-item :label="t('songRequests.settings.songs.minViews')">
					<n-input-number v-model:value="formValue.song!.minViews" :min="0" :max="99999999999999" />
				</n-form-item>
				<n-form-item :label="t('songRequests.settings.deniedSongs')">
					<n-select
						v-model:value="formValue.denyList!.songs"
						:loading="songsSearch.isLoading.value"
						remote
						filterable
						multiple
						:options="songsSearchOptions"
						:render-label="renderSelectOption"
						clearable
						@search="(v) => songsSearchValue = v"
					/>
				</n-form-item>
			</n-tab-pane>
		</n-tabs>

		<n-button secondary block type="success" @click="save">
			{{ t('sharedButtons.save') }}
		</n-button>
	</n-form>
</template>
