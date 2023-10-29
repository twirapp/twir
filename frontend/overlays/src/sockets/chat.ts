import { useWebSocket } from '@vueuse/core';
import { reactive, watch } from 'vue';

type Event = {
	eventName: string,
	data: Record<string, any>
	createdAt: string
}

type BadgeVersion = {
	id: string,
		image_url_1x: string,
		image_url_2x: string,
		image_url_4x: string,
}

type ChatBadge = {
	set_id: string,
	versions: Array<BadgeVersion>
}

export type Settings = {
	channelId: string,
	channelName: string,
	channelDisplayName: string,
	globalBadges: Map<string, ChatBadge>,
	channelBadges: Map<string, BadgeVersion>,
	messageTimeout: number,
}

export const useChatSocket = (apiKey: string) => {
	const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
	const host = window.location.host;
	const settings = reactive<Settings>({
		channelId: '',
		channelName: '',
		channelDisplayName: '',
		globalBadges: new Map<string, ChatBadge>(),
		channelBadges: new Map<string, BadgeVersion>(),
		messageTimeout: 0,
	});

	const { data, send } = useWebSocket(
		`${protocol}://${host}/socket/chat?apiKey=${apiKey}`,
		{
			immediate: true,
			autoReconnect: {
				delay: 500,
			},
			onConnected() {
				send(JSON.stringify({ eventName: 'getSettings' }));
			},
		},
	);

	watch(data, (d: string) => {
		const event = JSON.parse(d) as Event;

		if (event.eventName === 'settings') {
			for (const badge of event.data.globalBadges) {
				settings.globalBadges.set(badge.set_id, badge);
			}

			for (const badge of event.data.channelBadges) {
				for (const version of badge.versions) {
					settings.channelBadges.set(
						`${badge.set_id}-${version.id}`,
						version,
					);
				}
			}

			settings.channelId = event.data.channelId;
			settings.channelName = event.data.channelName;
			settings.channelDisplayName = event.data.channelDisplayName;
			settings.messageTimeout = event.data.messageTimeout ?? 0;
		}
	});

	return {
		settings,
	};
};
