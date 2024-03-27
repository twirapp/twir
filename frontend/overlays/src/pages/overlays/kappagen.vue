<script setup lang="ts">
import KappagenOverlay from '@twirapp/kappagen';
import type { Emote, KappagenAnimations, KappagenMethods } from '@twirapp/kappagen/types';
import { storeToRefs } from 'pinia';
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useRoute } from 'vue-router';

import { useKappagenEmotesBuilder } from '@/composables/kappagen/use-kappagen-builder.js';
import { useKappagenIframe } from '@/composables/kappagen/use-kappagen-iframe.js';
import { useKappagenSettings } from '@/composables/kappagen/use-kappagen-settings.js';
import { useKappagenOverlaySocket } from '@/composables/kappagen/use-kappagen-socket.js';
import { useChatTmi, type ChatSettings, type ChatMessage } from '@/composables/tmi/use-chat-tmi.js';

const kappagen = ref<KappagenMethods>();
const route = useRoute();
const { kappagenSettings, overlaySettings } = storeToRefs(useKappagenSettings());

const playAnimation = (emotes: Emote[], animation: KappagenAnimations) => {
	if (!kappagen.value) return Promise.resolve();
	return kappagen.value.playAnimation(emotes, animation);
};

const showEmotes = (emotes: Emote[]) => {
	if (!kappagen.value) return;
	kappagen.value.showEmotes(emotes);
};

const emotesBuilder = useKappagenEmotesBuilder();

const socket = useKappagenOverlaySocket({
	playAnimation,
	showEmotes,
	emotesBuilder,
});

const iframe = useKappagenIframe({
	playAnimation,
	showEmotes,
	clear: () => {
		kappagen.value?.clear();
	},
});

function onMessage(msg: ChatMessage): void {
	if (msg.type === 'system' || !overlaySettings.value?.enableSpawn) return;

	const firstChunk = msg.chunks.at(0)!;
	if (firstChunk.type === 'text' && firstChunk.value.startsWith('!')) {
		return;
	}

	const generatedEmotes = emotesBuilder.buildSpawnEmotes(msg.chunks);
	if (!generatedEmotes.length) return;
	showEmotes(generatedEmotes);
}

const chatSettings = computed<ChatSettings>(() => {
	return {
		channelId: overlaySettings.value?.channelId ?? '',
		channelName: overlaySettings.value?.channelName ?? '',
		emotes: {
			ffz: overlaySettings.value?.emotes?.ffzEnabled,
			bttv: overlaySettings.value?.emotes?.bttvEnabled,
			sevenTv: overlaySettings.value?.emotes?.sevenTvEnabled,
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
});

onUnmounted(() => {
	iframe.destroy();
	socket.destroy();
	destroy();
});
</script>

<template>
	<kappagen-overlay ref="kappagen" :config="kappagenSettings" :is-rave="overlaySettings?.enableRave" />
</template>
