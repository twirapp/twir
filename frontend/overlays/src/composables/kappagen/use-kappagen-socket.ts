import type { KappagenAnimations, KappagenMethods } from '@twirapp/kappagen/types';
import { useWebSocket } from '@vueuse/core';
import { storeToRefs } from 'pinia';
import { ref, watch } from 'vue';

import { type Buidler } from './use-kappagen-builder.js';

import type { TwirWebSocketEvent } from '@/api.js';
import { useKappagenSettings } from '@/composables/kappagen/use-kappagen-settings.js';
import { useMessageHelpers } from '@/composables/tmi/use-message-helpers.js';
import { generateSocketUrlWithParams } from '@/helpers.js';
import type { KappagenSettings, KappagenTriggerRequestEmote } from '@/types.js';

type Options = Omit<KappagenMethods, 'clear'> & {
	emotesBuilder: Buidler
}

export const useKappagenOverlaySocket = (options: Options) => {
	const { makeMessageChunks } = useMessageHelpers();
	const kappagenSettingsStore = useKappagenSettings();
	const { overlaySettings } = storeToRefs(kappagenSettingsStore);

	const kappagenUrl = ref('');
	const { data, send, open, close } = useWebSocket(
		kappagenUrl,
		{
			immediate: false,
			autoReconnect: {
				delay: 500,
			},
			onConnected() {
				send(JSON.stringify({ eventName: 'getSettings' }));
			},
		},
	);

	function randomAnimation() {
		if (!overlaySettings.value) return;
		const enabledAnimations = overlaySettings.value.animations
			.filter((animation) => animation.enabled);

		const index = Math.floor(Math.random() * enabledAnimations.length)
		return enabledAnimations[index] as KappagenAnimations;
	}

	watch(data, (d: string) => {
		const event = JSON.parse(d) as TwirWebSocketEvent;

		if (event.eventName === 'settings') {
			const data = event.data as KappagenSettings;
			kappagenSettingsStore.updateSettings(data);
		}

		if (event.eventName === 'event') {
			const generatedEmotes = options.emotesBuilder.buildKappagenEmotes([]);

			const animation = randomAnimation();
			if (!animation) return;

			options.playAnimation(generatedEmotes, animation);
		}

		if (event.eventName === 'kappagen') {
			const data = event.data as { text: string, emotes?: KappagenTriggerRequestEmote[] };

			const emotesList: Record<string, string[]> = {};
			if (data.emotes) {
				for (const emote of data.emotes) {
					emotesList[emote.id] = emote.positions;
				}
			}

			const chunks = makeMessageChunks(
				data.text,
				{
					isSmaller: false,
					emotesList,
				},
			);
			const emotesForKappagen = options.emotesBuilder.buildKappagenEmotes(chunks);

			const animation = randomAnimation();
			if (!animation) return;

			options.playAnimation(emotesForKappagen, animation);
		}
	});

	function destroy() {
		close();
	}

	function connect(apiKey: string): void {
		const url = generateSocketUrlWithParams('/overlays/kappagen', {
			apiKey,
		});

		kappagenUrl.value = url;

		open();
	}

	return {
		connect,
		destroy,
	};
};
