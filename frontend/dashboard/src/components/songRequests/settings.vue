<script setup lang='ts'>
import { useDebounce } from '@vueuse/core';
import {
	type SelectOption,
	NTabs,
	NTabPane,
	NSpace,
	NSwitch,
	NText,
	NInput,
	NInputNumber,
	NForm,
	NFormItem,
	NSelect,
	NAvatar,
} from 'naive-ui';
import { ref, computed, VNodeChild, h } from 'vue';

import { useTwitchRewards, useYoutubeVideoOrChannelSearch, YoutubeSearchType } from '@/api/index.js';

const formValue = ref({
	enabled: true,
	acceptOnlyWhenOnline: true,
	channelPointsRewardId: '',
	maxRequests: 500,
	announcePlay: true,
	neededVotesVorSkip: 30,
	denyList: {
		artistsNames: [],
		songs: [],
		users: [],
		channels: [],
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

const twitchRewards = useTwitchRewards();
const rewardsOptions = computed(() => {
	return twitchRewards.data?.value?.rewards.map((reward) => {
		return {
			label: reward.title,
			value: reward.id,
			image: reward.image?.url4X ?? '',
			disabled: !reward.isUserInputRequired,
		};
	}) ?? [];
});

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
const channelsSearch = useYoutubeVideoOrChannelSearch(
	channelsSearchDebounced,
	YoutubeSearchType.Channel,
);
const channelsSearchOptions = computed(() => {
	return channelsSearch.data?.value?.items.map((channel) => {
		return {
			label: channel.title,
			value: channel.id,
			image: channel.thumbnail,
		};
	}) ?? [];
});
</script>

<template>
  <n-form>
    <n-tabs
      type="line"
      animated
    >
      <n-tab-pane name="general" tab="General">
        <n-space vertical>
          <n-space justify="space-between">
            <n-text>Enabled</n-text>
            <n-switch v-model:value="formValue.enabled" />
          </n-space>

          <n-space justify="space-between">
            <n-text>Accept requests only when stream online</n-text>
            <n-switch v-model:value="formValue.acceptOnlyWhenOnline" />
          </n-space>

          <n-space justify="space-between">
            <n-text>Announce track play</n-text>
            <n-switch v-model:value="formValue.announcePlay" />
          </n-space>

          <n-form-item label="Needed percentage of votes for skip song" path="neededVotesVorSkip">
            <n-input-number v-model:value="formValue.neededVotesVorSkip" />
          </n-form-item>

          <n-form-item label="Channel reward for requesting songs" path="channelPointsRewardId">
            <n-select
              v-model:value="formValue.channelPointsRewardId"
              :loading="twitchRewards.isLoading.value"
              remote
              filterable
              :options="rewardsOptions"
              :render-label="renderSelectOption as any"
              clearable
            />
          </n-form-item>

          <n-form-item label="Denied channels" path="channelPointsRewardId">
            <n-select
              v-model:value="formValue.denyList.channels"
              :loading="channelsSearch.isLoading.value"
              remote
              filterable
              :options="channelsSearchOptions"
              clearable
              :render-label="renderSelectOption as any"
              @search="(v) => channelsSearchValue = v"
            />
          </n-form-item>
        </n-space>
      </n-tab-pane>

      <n-tab-pane name="users" tab="Users">
        s
      </n-tab-pane>

      <n-tab-pane name="songs" tab="Songs">
        t
      </n-tab-pane>

      <n-tab-pane name="translations" tab="Translations">
        t
      </n-tab-pane>
    </n-tabs>
  </n-form>
</template>

