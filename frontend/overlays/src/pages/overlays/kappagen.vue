<script setup lang="ts">
import KappagenOverlay from '@twirapp/kappagen';
import type { KappagenEmoteConfig } from '@twirapp/kappagen';
import { storeToRefs } from 'pinia';
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import '@twirapp/kappagen/style.css';

import { useKappagenEmotesBuilder } from '@/composables/kappagen/use-kappagen-builder.js';
import { useKappagenIframe } from '@/composables/kappagen/use-kappagen-iframe.js';
import { useKappagenSettings } from '@/composables/kappagen/use-kappagen-settings.js';
import { useKappagenOverlaySocket } from '@/composables/kappagen/use-kappagen-socket.js';
import { useChatTmi, type ChatSettings, type ChatMessage } from '@/composables/tmi/use-chat-tmi.js';
import type { KappagenSpawnAnimatedEmotesFn, KappagenSpawnEmotesFn } from '@/types.js';

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

const kappagenCallback: KappagenSpawnAnimatedEmotesFn = (emotes, animation) => {
	kappagen.value?.kappagen.run(emotes, animation);
};

const spawnCallback: KappagenSpawnEmotesFn = (emotes) => {
	kappagen.value?.emote.addEmotes(emotes);
	kappagen.value?.emote.showEmotes();
};

const emojiStyle = computed(() => settings.value?.emotes?.emojiStyle);
const emotesBuilder = useKappagenEmotesBuilder(emojiStyle);

const socket = useKappagenOverlaySocket({
	kappagenCallback,
	spawnCallback,
	emotesBuilder,
});

const iframe = useKappagenIframe({
	kappagenCallback,
	spawnCallback,
	clearCallback: () => {
		kappagen.value?.clear();
	},
});

function onMessage(msg: ChatMessage): void {
	if (msg.type === 'system' || !settings.value?.enableSpawn) return;

	const firstChunk = msg.chunks.at(0)!;
	if (firstChunk.type === 'text' && firstChunk.value.startsWith('!')) {
		return;
	}

	const generatedEmotes = emotesBuilder.buildSpawnEmotes(msg.chunks);
	if (!generatedEmotes.length) return;
	kappagen.value?.emote.addEmotes(generatedEmotes);
	kappagen.value?.emote.showEmotes();
}

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: settings.value?.channelId ?? '',
		channelName: settings.value?.channelName ?? '',
		emotes: {
			ffz: settings.value?.emotes?.ffzEnabled,
			bttv: settings.value?.emotes?.bttvEnabled,
			sevenTv: settings.value?.emotes?.sevenTvEnabled,
		},
		onMessage,
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
