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
		const event = JSON.parse(d) as TwirWebSocketEvent;

		if (event.eventName === 'settings') {
			if (overlayId.value && event.data.id !== overlayId.value) return;

			const data = event.data;

			for (const badge of data.global_badges) {
				settings.value.globalBadges.set(badge.set_id, badge);
			}

			for (const badge of data.channel_badges) {
				for (const version of badge.versions) {
					settings.value.channelBadges.set(
						`${badge.set_id}-${version.id}`,
						version,
					);
				}
			}

			settings.value.channelId = data.channel_id;
			settings.value.channelName = data.channel_name;
			settings.value.channelDisplayName = data.channel_display_name;
			settings.value.messageHideTimeout = data.message_hide_timeout;
			settings.value.messageShowDelay = data.message_show_delay;
			settings.value.preset = data.preset;
			settings.value.fontSize = data.font_size;
			settings.value.hideBots = data.hide_bots;
			settings.value.hideCommands = data.hide_commands;
			settings.value.fontFamily = data.font_family;
			settings.value.showAnnounceBadge = data.show_announce_badge;
			settings.value.showBadges = data.show_badges;
			settings.value.textShadowColor = data.text_shadow_color;
			settings.value.textShadowSize = data.text_shadow_size;
			settings.value.chatBackgroundColor = data.chat_background_color;
			settings.value.direction = data.direction;
			settings.value.fontStyle = data.font_style;
			settings.value.fontWeight = data.font_weight;
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
