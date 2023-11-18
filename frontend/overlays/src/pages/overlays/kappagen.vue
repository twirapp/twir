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
import { useThirdPartyEmotes } from '../../components/chat_tmi_emotes.js';
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
		emoteConfig.size.min = s.size.min;
		emoteConfig.size.max = s.size.max;
		emoteConfig.size.ratio!.normal = 1/(s.size.ratioNormal ?? 12);
		emoteConfig.size.ratio!.small = 1/(s.size.ratioSmall ?? 24);
	}

	if (s.cube) {
		emoteConfig.cube.speed = s.cube.speed;
	}

	if (s.animation) {
		emoteConfig.in.fade = s.animation.fadeIn;
		emoteConfig.in.zoom = s.animation.zoomIn;

		emoteConfig.out.fade = s.animation.fadeOut;
		emoteConfig.out.zoom = s.animation.zoomOut;
	}

	if (typeof s.enableRave !== 'undefined') {
		emoteConfig.rave = s.enableRave;
	}
});

const emotesBuilder = useKappagenEmotesBuilder();

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

const socket = useKappagenOverlaySocket(apiKey, {
	kappagenCallback,
	setSettingsCallback,
	spawnCallback,
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
	};
});

const channelId = computed(() => kappagenSettings.value?.channelId ?? '');
const channelName = computed(() => kappagenSettings.value?.channelName ?? '');
useThirdPartyEmotes(channelName, channelId);

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: channelId.value,
		channelName: channelName.value,
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
./kappagen/kappagen_overlay.js./kappagen/kappagen.js
./kappagen/settingsStore.js
