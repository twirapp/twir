<script setup lang="ts">
import KappagenOverlay from 'kappagen';
import type { KappagenEmoteConfig } from 'kappagen';
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import 'kappagen/style.css';

import { useKappagenBuilder as useKappagenEmotesBuilder } from './kappagen/builder.js';
import { useIframe } from './kappagen/iframe.js';
import { useChannelSettings } from './kappagen/settingsStore.js';
import { useKappagenOverlaySocket } from './kappagen/socket.js';
import type { KappagenCallback, SetSettingsCallback, SpawnCallback } from './kappagen/types.js';
import { useThirdPartyEmotes, type Opts as EmotesOpts } from '../../components/chat_tmi_emotes.js';
import { ChatSettings, useTmiChat } from '../../sockets/chat_tmi';

const kappagen = ref<InstanceType<typeof KappagenOverlay>>();
const route = useRoute();
const apiKey = route.params.apiKey as string;

type Config = KappagenEmoteConfig & { rave: boolean }

const emoteConfig = reactive<Required<Config>>({
  max: 0,
  time: 5,
  queue: 0,
  cube: {
    speed: 6,
  },
	animation: {
		fade: {
			in: 8,
			out: 9,
		},
		zoom: {
			in: 17,
			out: 8,
		},
	},
	in: {
		fade: true,
		zoom: true,
	},
	out: {
		fade: true,
		zoom: true,
	},
	size: {
		min: 1,
		max: 256,
		ratio: {
			normal: 1/12,
			small: 1/24,
		},
	},
	rave: false,
});

const { kappagenSettings, setSettings } = useChannelSettings();
watch(kappagenSettings, (s) => {
	if (!s) return;

	if (s.emotes) {
		emoteConfig.max = s.emotes.max;
		emoteConfig.time = s.emotes.time;
		emoteConfig.queue = s.emotes.queue;
	}

	if (s.size) {
		emoteConfig.size = {
			min: s.size.min,
			max: s.size.max,
			ratio: {
				normal: s.size.ratioNormal,
				small: s.size.ratioSmall,
			},
		};
	}

	if (s.cube) {
		emoteConfig.cube = {
			speed: s.cube.speed,
		};
	}

	if (s.animation) {
		emoteConfig.in = {
			fade: s.animation.fadeIn,
			zoom: s.animation.zoomIn,
		};

		emoteConfig.out = {
			fade: s.animation.fadeOut,
			zoom: s.animation.zoomOut,
		};
	}

	if (typeof s.enableRave !== 'undefined') {
		emoteConfig.rave = s.enableRave;
	}
});

const kappagenCallback: KappagenCallback = (emotes, animation) => {
	kappagen.value?.kappagen.run(emotes, animation);
};

const spawnCallback: SpawnCallback = (emotes) => {
	kappagen.value?.emote.addEmotes(emotes);
	kappagen.value?.emote.showEmotes();
};

const setSettingsCallback: SetSettingsCallback = (settings) => {
	setSettings(settings);
};


const emojiStyle = computed(() => kappagenSettings.value?.emotes?.emojiStyle);
const emotesBuilder = useKappagenEmotesBuilder(emojiStyle);

const socket = useKappagenOverlaySocket(apiKey, {
	kappagenCallback,
	setSettingsCallback,
	spawnCallback,
	emotesBuilder,
});
const iframe = useIframe({
	kappagenCallback,
	setSettingsCallback,
	spawnCallback,
	clearCallback: () => {
		kappagen.value?.clear();
	},
});

onMounted(() => {
	if (window.frameElement) {
		iframe.create();
	} else {
		socket.create();
	}

	kappagen.value?.init();

	return () => {
		iframe.destroy();
		socket.destroy();
	};
});

const emotesOpts = computed<EmotesOpts>(() => {
	return {
		channelName: kappagenSettings.value?.channelName,
		channelId: kappagenSettings.value?.channelId,
		ffz: kappagenSettings.value?.emotes?.ffzEnabled,
		bttv: kappagenSettings.value?.emotes?.bttvEnabled,
		sevenTv: kappagenSettings.value?.emotes?.sevenTvEnabled,
	};
});

useThirdPartyEmotes(emotesOpts);

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: emotesOpts.value.channelId!,
		channelName: emotesOpts.value.channelName!,
		onMessage: (msg) => {
			if (msg.type === 'system') return;

			const firstChunk = msg.chunks.at(0)!;
			if (firstChunk.type === 'text' && firstChunk.value.startsWith('!')) {
				return;
			}

			if (!kappagenSettings.value?.enableSpawn) return;

			const generatedEmotes = emotesBuilder.buildSpawnEmotes(msg.chunks);
			if (!generatedEmotes.length) return;
			kappagen.value?.emote.addEmotes(generatedEmotes);
			kappagen.value?.emote.showEmotes();
		},
	};
});

const { destroy } = useTmiChat(chatSettings);

onUnmounted(() => {
	destroy();
});
</script>

<template>
	<kappagen-overlay ref="kappagen" :emote-config="emoteConfig" :is-rave="emoteConfig.rave" />
</template>
