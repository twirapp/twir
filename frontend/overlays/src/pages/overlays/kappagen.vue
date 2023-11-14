<script setup lang="ts">
import type { Settings } from '@twir/grpc/generated/api/api/overlays_kappagen';
import type { TriggerKappagenRequest_Emote } from '@twir/grpc/generated/websockets/websockets';
import KappagenOverlay, { KappagenAnimations, type KappagenEmoteConfig } from 'kappagen';
import { computed, onMounted, onUnmounted, reactive, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import 'kappagen/style.css';

import { useThirdPartyEmotes } from '../../components/chat_tmi_emotes.js';
import { makeMessageChunks } from '../../components/chat_tmi_helpers';
import { useKappagenBuilder, twirEmote } from '../../components/kappagen';
import { animations } from '../../components/kappagen_animations.js';
import { ChatSettings, useTmiChat } from '../../sockets/chat_tmi';
import { useKappagenOverlaySocket } from '../../sockets/kappagen_overlay';
import { TwirWebSocketEvent } from '../../sockets/types';

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
const channelSettings = ref<IncomingSettings>();

const socket = useKappagenOverlaySocket(apiKey);

const builder = useKappagenBuilder();

const onWindowMessage = (msg: MessageEvent<string>) => {
	const parsedData = JSON.parse(msg.data) as { key: string, data?: any };
	kappagen.value?.clear();

	if (parsedData.key === 'settings' && parsedData.data) {
		const settings = parsedData.data as Settings & { channelName: string, channelId: string };
		setSettings(settings);

		kappagen.value?.kappagen.run(
			[twirEmote],
			animations[Math.floor(Math.random() * animations.length)],
			);
		}

	if (parsedData.key === 'kappa') {
		kappagen.value?.kappagen.run(
			[twirEmote],
			animations[Math.floor(Math.random() * animations.length)],
		);
	}

	if (parsedData.key === 'kappaWithAnimation') {
		kappagen.value?.kappagen.run(
			[twirEmote],
			parsedData.data.animation,
		);
	}

	if (parsedData.key === 'spawn') {
		kappagen.value?.emote.addEmotes([twirEmote]);
		kappagen.value?.emote.showEmotes();
	}
};

type IncomingSettings = Settings & { channelId: string, channelName: string, kappagenCommandName?: string }


const channelId = computed(() => channelSettings.value?.channelId ?? '');
const channelName = computed(() => channelSettings.value?.channelName ?? '');
const { emotes } = useThirdPartyEmotes(channelName, channelId);

const setSettings = (settings: IncomingSettings) => {
	if (settings.emotes) {
		emoteConfig.max = settings.emotes.max;
		emoteConfig.time = settings.emotes.time;
		emoteConfig.queue = settings.emotes.queue;
	}

	if (settings.size) {
		emoteConfig.size.min = settings.size.min;
		emoteConfig.size.max = settings.size.max;
		emoteConfig.size.ratio!.normal = 1/(settings.size.ratioNormal ?? 12);
		emoteConfig.size.ratio!.small = 1/(settings.size.ratioSmall ?? 24);
	}

	if (settings.cube) {
		emoteConfig.cube.speed = settings.cube.speed;
	}

	if (settings.animation) {
		emoteConfig.in.fade = settings.animation.fadeIn;
		emoteConfig.in.zoom = settings.animation.zoomIn;

		emoteConfig.out.fade = settings.animation.fadeOut;
		emoteConfig.out.zoom = settings.animation.zoomOut;
	}

	if (typeof settings.enableRave !== 'undefined') {
		emoteConfig.rave = settings.enableRave;
	}
};

onMounted(() => {
	if (window.frameElement) {
		window.postMessage('getSettings');
		window.addEventListener('message', onWindowMessage);
	} else {
		socket.connect();
	}

	kappagen.value?.init();

	return () => {
		window.removeEventListener('message', onWindowMessage);
	};
});

watch(socket.data, (d: string) => {
	const event = JSON.parse(d) as TwirWebSocketEvent;

	if (event.eventName === 'settings') {
		const data = event.data as IncomingSettings;
		channelSettings.value = data;
		setSettings(data);
	}

	if (event.eventName === 'event') {
		if (!channelSettings.value) return;
		const generatedEmotes = Object.values(emotes.value)
			.filter(e => !e.isModifier && !e.isZeroWidth)
			.map(e => ({ url: e.urls.at(-1)! }));

		const enabledAnimations = channelSettings.value.animations.filter(a => a.enabled);
		const randomAnimation = enabledAnimations[Math.floor(Math.random()*enabledAnimations.length)];

		kappagen.value?.kappagen.run(generatedEmotes, randomAnimation as KappagenAnimations);
	}

	if (event.eventName === 'kappagen') {
		const data = event.data as { text: string, emotes: TriggerKappagenRequest_Emote[] };
		console.log(data)
		const emotes = builder.buildKappagenEmotes(makeMessageChunks(
			data.text,
			data.emotes.reduce((acc, curr) => {
				acc[curr.id] = curr.positions;
				return acc;
			}, {} as Record<string, string[]>),
		));

		if (!emotes.length || !channelSettings.value) return;

		const enabledAnimations = channelSettings.value.animations.filter(a => a.enabled);
		const randomAnimation = enabledAnimations[Math.floor(Math.random()*enabledAnimations.length)];
		kappagen.value?.kappagen.run(emotes, randomAnimation as KappagenAnimations);
	}
});

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

			const generatedEmotes = builder.buildSpawnEmotes(msg.chunks);
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
