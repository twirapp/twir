<script setup lang='ts'>
import type { YouTubeSettings } from '@twir/grpc/generated/api/api/modules_sr';
import { useDebounce } from '@vueuse/core';
import { NSpace, SelectOption, NAvatar, NText  } from 'naive-ui';
import { ref, computed, VNodeChild, h, watch, unref, toRaw, reactive } from 'vue';

import { useTwitchRewards, useYoutubeVideoOrChannelSearch, YoutubeSearchType } from '@/api/index.js';
import { useYoutubeModuleSettings } from '@/api/index.js';
import TwitchSearchUsers from '@/components/twitchUsers/multiple.vue';

const youtubeModuleManager = useYoutubeModuleSettings();
const youtubeModuleData = youtubeModuleManager.getAll();
const youtubeModuleSettings = computed(() => {
	return youtubeModuleData.data?.value?.data ?? null;
});
const youtubeModuleUpdater = youtubeModuleManager.update();

const formValue = ref<YouTubeSettings>({
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

watch(youtubeModuleSettings, async (v) => {
	if (!v) return;

	formValue.value = v;
});

async function save() {
	const data = unref(formValue);

	await youtubeModuleUpdater.mutateAsync({ data });
}

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
              :disabled="twitchRewards.isError.value"
            />
          </n-form-item>

          <n-form-item label="Denied channels" path="channelPointsRewardId">
            {{ formValue.denyList!.channels }}
            <n-select
              v-model:value="formValue.denyList!.channels"
              :loading="selectedChannels.isLoading.value"
              remote
              filterable
              :options="channelsOptions"
              clearable
              multiple
              :clear-filter-after-select="false"
              :render-label="renderSelectOption as any"
              @search="(v) => channelsSearchValue = v"
            />
          </n-form-item>
        </n-space>
      </n-tab-pane>

      <n-tab-pane name="users" tab="Users">
        <n-form-item label="Maximum songs by user in queue" path="user.maxRequests">
          <n-input-number v-model:value="formValue.user!.maxRequests" :min="0" :max="1000" />
        </n-form-item>
        <n-form-item
          label="Minimal watch time of user for request song (minutes)"
          path="user.minWatchTime"
        >
          <n-input-number v-model:value="formValue.user!.minWatchTime" :min="0" :max="999999999" />
        </n-form-item>
        <n-form-item
          label="Minimal messages by user for request song"
          path="user.minMessages"
        >
          <n-input-number v-model:value="formValue.user!.minMessages" :min="0" :max="999999999" />
        </n-form-item>
        <n-form-item
          label="Minimal follow time for request song (minutes)"
          path="user.minFollowTime"
        >
          <n-input-number v-model:value="formValue.user!.minFollowTime" :min="0" :max="99999999999999" />
        </n-form-item>

        <n-form-item label="Denied users for requests">
          <twitch-search-users v-model="formValue.denyList!.users" />
        </n-form-item>
      </n-tab-pane>

      <n-tab-pane name="songs" tab="Songs">
        <n-form-item label="Maximum number of songs in queue">
          <n-input-number v-model:value="formValue.maxRequests" :min="0" :max="99999999999999" />
        </n-form-item>
        <n-form-item label="Min length of song for request (minutes)">
          <n-input-number v-model:value="formValue.song!.minLength" :min="0" :max="99999999999999" />
        </n-form-item>
        <n-form-item label="Max length of song for request (minutes)">
          <n-input-number v-model:value="formValue.song!.maxLength" :min="0" :max="99999999999999" />
        </n-form-item>
        <n-form-item label="Minimal views on song for request">
          <n-input-number v-model:value="formValue.song!.minViews" :min="0" :max="99999999999999" />
        </n-form-item>
        <n-form-item label="Denied songs for request">
          <n-select
            v-model:value="formValue.denyList!.songs"
            :loading="songsSearch.isLoading.value"
            remote
            filterable
            :options="songsSearchOptions"
            :render-label="renderSelectOption as any"
            clearable
            @search="(v) => songsSearchValue = v"
          />
        </n-form-item>
      </n-tab-pane>

      <n-tab-pane name="translations" tab="Translations">
        <div v-for="(key) in Object.keys(formValue.translations)" :key="key">
          <n-form-item v-if="typeof formValue.translations?.[key] === 'string'" :label="key" :path="key">
            <n-input
              v-model:value="formValue.translations[key]"
              type="textarea"
              :autosize="{ minRows: 1, maxRows: 6 }"
            />
          </n-form-item>

          <div v-else>
            <n-form-item
              v-for="subKey of Object.keys(formValue.translations?.[key])"
              :key="subKey"
              :label="subKey"
              :path="subKey"
            >
              <n-input
                v-model:value="formValue.translations![key][subKey]"
                type="textarea"
                :autosize="{ minRows: 1, maxRows: 6 }"
              />
            </n-form-item>
          </div>
        </div>
      </n-tab-pane>
    </n-tabs>

    <n-button secondary block type="success" @click="save">
      Save
    </n-button>
  </n-form>
</template>

