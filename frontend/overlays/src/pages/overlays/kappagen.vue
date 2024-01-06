<script setup lang="ts">
import KappagenOverlay from 'kappagen';
import type { KappagenEmoteConfig } from 'kappagen';
import { storeToRefs } from 'pinia';
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import 'kappagen/style.css';

import {
	useKappagenBuilder as useKappagenEmotesBuilder,
} from '@/composables/kappagen/use-kappagen-builder.js';
import { useKappagenIframe } from '@/composables/kappagen/use-kappagen-iframe.js';
import { useKappagenSettings } from '@/composables/kappagen/use-kappagen-settings.js';
import { useKappagenOverlaySocket } from '@/composables/kappagen/use-kappagen-socket.js';
import { useChatTmi, type ChatSettings } from '@/composables/tmi/use-chat-tmi.ts';
import {
	useThirdPartyEmotes,
	type ThirdPartyEmotesOptions,
} from '@/composables/tmi/use-third-party-emotes.ts';
import type { KappagenCallback, SetSettingsCallback, SpawnCallback } from '@/types.js';

const kappagen = ref<InstanceType<typeof KappagenOverlay>>();
const route = useRoute();

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
			normal: 1 / 12,
			small: 1 / 24,
		},
	},
	rave: false,
});

const kappagenSettingsStore = useKappagenSettings();
const { settings } = storeToRefs(kappagenSettingsStore);

watch(() => settings.value, (settings) => {
	if (!settings) return;

	if (settings.emotes) {
		emoteConfig.max = settings.emotes.max;
		emoteConfig.time = settings.emotes.time;
		emoteConfig.queue = settings.emotes.queue;
	}

	if (settings.size) {
		emoteConfig.size = {
			min: settings.size.min,
			max: settings.size.max,
			ratio: {
				normal: settings.size.ratioNormal,
				small: settings.size.ratioSmall,
			},
		};
	}

	if (settings.cube) {
		emoteConfig.cube = {
			speed: settings.cube.speed,
		};
	}

	if (settings.animation) {
		emoteConfig.in = {
			fade: settings.animation.fadeIn,
			zoom: settings.animation.zoomIn,
		};

		emoteConfig.out = {
			fade: settings.animation.fadeOut,
			zoom: settings.animation.zoomOut,
		};
	}

	if (typeof settings.enableRave !== 'undefined') {
		emoteConfig.rave = settings.enableRave;
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
	kappagenSettingsStore.setSettings(settings);
};

const emojiStyle = computed(() => settings.value?.emotes?.emojiStyle);
const emotesBuilder = useKappagenEmotesBuilder(emojiStyle);

const socket = useKappagenOverlaySocket({
	kappagenCallback,
	setSettingsCallback,
	spawnCallback,
	emotesBuilder,
});

const iframe = useKappagenIframe({
	kappagenCallback,
	setSettingsCallback,
	spawnCallback,
	clearCallback: () => {
		kappagen.value?.clear();
	},
});

const emotesOptions = computed<ThirdPartyEmotesOptions>(() => {
	return {
		channelName: settings.value?.channelName,
		channelId: settings.value?.channelId,
		ffz: settings.value?.emotes?.ffzEnabled,
		bttv: settings.value?.emotes?.bttvEnabled,
		sevenTv: settings.value?.emotes?.sevenTvEnabled,
	};
});

useThirdPartyEmotes(emotesOptions);

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: emotesOptions.value.channelId!,
		channelName: emotesOptions.value.channelName!,
		onMessage: (msg) => {
			if (msg.type === 'system') return;

			const firstChunk = msg.chunks.at(0)!;
			if (firstChunk.type === 'text' && firstChunk.value.startsWith('!')) {
				return;
			}

			if (!settings.value?.enableSpawn) return;

			const generatedEmotes = emotesBuilder.buildSpawnEmotes(msg.chunks);
			if (!generatedEmotes.length) return;
			kappagen.value?.emote.addEmotes(generatedEmotes);
			kappagen.value?.emote.showEmotes();
		},
	};
});

const { destroy } = useChatTmi(chatSettings);

onMounted(() => {
	if (window.frameElement) {
		iframe.create();
	} else {
		const apiKey = route.params.apiKey as string;
		socket.connect(apiKey);
	}

	kappagen.value?.init();
});

onUnmounted(() => {
	iframe.destroy();
	socket.destroy();
	destroy();
});
</script>

<template>
	<kappagen-overlay ref="kappagen" :emote-config="emoteConfig" :is-rave="emoteConfig.rave" />
</template>
