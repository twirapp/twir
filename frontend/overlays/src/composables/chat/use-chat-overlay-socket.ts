import type { Settings, ChatBadge, BadgeVersion } from '@twir/frontend-chat';
import { useWebSocket } from '@vueuse/core';
import { defineStore } from 'pinia';
import { ref, watch } from 'vue';

import type { TwirWebSocketEvent } from '@/api.js';
import { generateSocketUrlWithParams } from '@/helpers.js';

export const useChatOverlaySocket = defineStore('chat-socket', () => {
	const settings = ref<Settings>({
		channelId: '',
		channelName: '',
		channelDisplayName: '',
		globalBadges: new Map<string, ChatBadge>(),
		channelBadges: new Map<string, BadgeVersion>(),
		messageHideTimeout: 0,
		messageShowDelay: 0,
		preset: 'clean',
		fontSize: 20,
		hideBots: false,
		hideCommands: false,
		fontFamily: 'Roboto',
		showAnnounceBadge: true,
		showBadges: true,
		textShadowColor: '',
		textShadowSize: 0,
		chatBackgroundColor: '',
		direction: 'top',
		fontStyle: 'normal',
		fontWeight: 400,
		paddingContainer: 0,
	});

	const overlayId = ref<string | undefined>();
	const socketUrl = ref('');

	const { data, status, send, open, close } = useWebSocket(
		socketUrl,
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
		const event = JSON.parse(d) as TwirWebSocketEvent<Settings>;
		if (event.eventName === 'settings') {
			if (overlayId.value && event.data.id !== overlayId.value) return;

			const data = event.data;

			for (const [, badge] of data.globalBadges) {
				settings.value.globalBadges.set(badge.set_id, badge);
			}

			// :)
			for (const [, badge] of data.channelBadges as any as Map<string, ChatBadge>) {
				for (const version of badge.versions) {
					settings.value.channelBadges.set(
						`${badge.set_id}-${version.id}`,
						version,
					);
				}
			}

			settings.value = data;
		}
	});

	function destroy(): void {
		close();
	}

	function connect(apiKey: string, _overlayId?: string): void {
		if (status.value === 'OPEN') return;

		const url = generateSocketUrlWithParams('/overlays/chat', {
			apiKey,
			id: _overlayId,
		});

		socketUrl.value = url;
		overlayId.value = _overlayId;

		open();
	}

	return {
		settings,
		connect,
		destroy,
	};
});
