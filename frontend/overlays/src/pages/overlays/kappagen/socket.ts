import type { TriggerKappagenRequest_Emote } from '@twir/grpc/generated/websockets/websockets';
import { useWebSocket } from '@vueuse/core';
import type { Emote, KappagenAnimations } from 'kappagen';
import { watch } from 'vue';

import { useKappagenBuilder } from './builder.js';
import { kappagenSettings } from './settingsStore.js';
import type { KappagenCallback, SpawnCallback, SetSettingsCallback, KappagenSettings } from './types.js';
import { emotes } from '../../../components/chat_tmi_emotes.js';
import { makeMessageChunks } from '../../../components/chat_tmi_helpers.js';
import type { TwirWebSocketEvent } from '../../../sockets/types';

type Opts = {
	kappagenCallback: KappagenCallback
	spawnCallback: SpawnCallback
	setSettingsCallback: SetSettingsCallback
}

export const useKappagenOverlaySocket = (apiKey: string, opts: Opts) => {
	const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
	const host = window.location.host;
	const emotesBuilder = useKappagenBuilder();

	const { data, send, open, close } = useWebSocket(
		`${protocol}://${host}/socket/overlays/kappagen?apiKey=${apiKey}`,
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


	watch(data, (d: string) => {
		const event = JSON.parse(d) as TwirWebSocketEvent;

		if (event.eventName === 'settings') {
			const data = event.data as KappagenSettings;
			opts.setSettingsCallback(data);
		}

		if (event.eventName === 'event') {
			if (!kappagenSettings.value) return;

			const generatedEmotes: Emote[] = Object.values(emotes.value)
				.filter(e => !e.isModifier && !e.isZeroWidth)
				.map(e => ({ url: e.urls.at(-1)!, width: e.width, height: e.height }));

			const enabledAnimations = kappagenSettings.value.animations.filter(a => a.enabled);
			const randomAnimation = enabledAnimations[Math.floor(Math.random()*enabledAnimations.length)];

			opts.kappagenCallback(generatedEmotes, randomAnimation as KappagenAnimations);
		}

		if (event.eventName === 'kappagen') {
			const data = event.data as { text: string, emotes?: TriggerKappagenRequest_Emote[] };
			const emotes = emotesBuilder.buildKappagenEmotes(makeMessageChunks(
				data.text,
				data.emotes?.reduce((acc, curr) => {
					acc[curr.id] = curr.positions;
					return acc;
				}, {} as Record<string, string[]>),
			));

			if (!emotes.length || !kappagenSettings.value) return;

			const enabledAnimations = kappagenSettings.value.animations.filter(a => a.enabled);
			const randomAnimation = enabledAnimations[Math.floor(Math.random()*enabledAnimations.length)];

			opts.kappagenCallback(emotes, randomAnimation as KappagenAnimations);
		}
	});

	const create = () => {
		open();
	};

	const destroy = () => {
		close();
	};

	return {
		create,
		destroy,
	};
};
