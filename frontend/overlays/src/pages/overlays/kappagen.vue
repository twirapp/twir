<script setup lang="ts">
import type { Settings } from '@twir/grpc/generated/api/api/overlays_kappagen';
import KappagenOverlay, { type KappagenEmoteConfig } from 'kappagen';
import { onMounted, reactive, ref, toRef } from 'vue';
import { useRoute } from 'vue-router';
import 'kappagen/style.css';

import { useThirdPartyEmotes } from '../../components/chat_tmi_emotes.js';
import { useKappagenBuilder, twirEmote } from '../../components/kappagen';
import { animations } from '../../components/kappagen_animations.js';
import { useKappagenOverlaySocket } from '../../sockets/kappagen_overlay';

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

useThirdPartyEmotes(toRef('fukushine'), toRef('971211575'));
const socket = useKappagenOverlaySocket(apiKey);

const builder = useKappagenBuilder();

const onWindowMessage = (msg: MessageEvent<string>) => {
	const parsedData = JSON.parse(msg.data) as { key: string, data?: any };

	if (parsedData.key === 'settings' && parsedData.data) {
		const settings = parsedData.data as Settings;
		setSettings(settings);
		kappagen.value?.clear();

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

	if (parsedData.key === 'spawn') {
		kappagen.value?.emote.addEmotes([twirEmote]);
		kappagen.value?.emote.showEmotes();
	}
};

const setSettings = (settings: Settings) => {
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
	if (window.parent) {
		window.addEventListener('message', onWindowMessage);
	} else {
		socket.connect();
	}

	kappagen.value?.init();

	return () => {
		window.removeEventListener('message', onWindowMessage);
	};
});
</script>

<template>
	<kappagen-overlay ref="kappagen" :emote-config="emoteConfig" :is-rave="emoteConfig.rave" />
</template>
