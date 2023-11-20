import type { TriggerKappagenRequest_Emote } from '@twir/grpc/generated/websockets/websockets';
import { useWebSocket } from '@vueuse/core';
import type { KappagenAnimations } from 'kappagen';
import { watch } from 'vue';

import { useKappagenBuilder } from './builder.js';
import { kappagenSettings } from './settingsStore.js';
import type { KappagenCallback, SpawnCallback, SetSettingsCallback, KappagenSettings } from './types.js';
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


	const randomAnimation = () => {
		if (!kappagenSettings.value) return;
		const enabledAnimations = kappagenSettings.value.animations.filter(a => a.enabled);
		return enabledAnimations[Math.floor(Math.random()*enabledAnimations.length)];
	};

	watch(data, (d: string) => {
		const event = JSON.parse(d) as TwirWebSocketEvent;

		if (event.eventName === 'settings') {
			const data = event.data as KappagenSettings;
			opts.setSettingsCallback(data);
		}

		if (event.eventName === 'event') {
			if (!kappagenSettings.value) return;

			const generatedEmotes = emotesBuilder.buildKappagenEmotes([]);

			const animation = randomAnimation();
			if (!animation) return;

			opts.kappagenCallback(generatedEmotes, animation as KappagenAnimations);
		}

		if (event.eventName === 'kappagen') {
			if (!kappagenSettings.value) return;

			const data = event.data as { text: string, emotes?: TriggerKappagenRequest_Emote[] };
			const emotesForKappagen = emotesBuilder.buildKappagenEmotes(makeMessageChunks(
				data.text,
				data.emotes?.reduce((acc, curr) => {
					acc[curr.id] = curr.positions;
					return acc;
				}, {} as Record<string, string[]>),
			));

			const animation = randomAnimation();
			if (!animation) return;

			opts.kappagenCallback(emotesForKappagen, animation as KappagenAnimations);
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
