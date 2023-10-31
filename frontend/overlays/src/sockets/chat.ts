import type { Settings, ChatBadge, BadgeVersion } from '@twir/frontend-chat';
import { useWebSocket } from '@vueuse/core';
import { type Ref, ref, watch } from 'vue';

type Event = {
	eventName: string,
	data: Record<string, any>
	createdAt: string
}

export const useChatSocket = (apiKey: string): { settings: Ref<Settings> } => {
	const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
	const host = window.location.host;
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
				settings.value.globalBadges.set(badge.set_id, badge);
			}

			for (const badge of event.data.channelBadges) {
				for (const version of badge.versions) {
					settings.value.channelBadges.set(
						`${badge.set_id}-${version.id}`,
						version,
					);
				}
			}

			const valuesForSet = Object.entries(event.data)
				.filter(([key]) => !['channelBadges', 'globalBadges'].includes(key));

			for (const [key, value] of valuesForSet) {
				// eslint-disable-next-line @typescript-eslint/ban-ts-comment
				// @ts-ignore
				settings.value[key] = value;
			}
		}
	});

	return {
		settings,
	};
};
