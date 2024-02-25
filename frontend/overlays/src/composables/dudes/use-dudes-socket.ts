import type { Settings } from '@twir/api/messages/overlays_dudes/overlays_dudes';
import type { DudesJumpRequest, DudesUserPunishedRequest } from '@twir/grpc/websockets/websockets';
import type {
	DudesGrowRequest,
	DudesChangeColorRequest,
	DudesUserSettings,
} from '@twir/types/overlays';
import { useWebSocket } from '@vueuse/core';
import { defineStore, storeToRefs } from 'pinia';
import { ref, watch } from 'vue';

import { useDudesSettings, type DudesConfig } from './use-dudes-settings';
import { useDudes } from './use-dudes.js';

import type { TwirWebSocketEvent } from '@/api.js';
import { generateSocketUrlWithParams } from '@/helpers.js';
import type { ChannelData } from '@/types.js';

const soundsDefaults: Partial<DudesConfig['dudes']['dude']['sounds']> = {
	enabled: false,
	volume: 0,
};

const nameBoxDefaults: Partial<Settings['nameBoxSettings']> = {
	strokeThickness: 0,
	fillGradientType: 0,
	dropShadow: false,
	dropShadowAlpha: 0,
	dropShadowBlur: 0,
	dropShadowDistance: 0,
	dropShadowAngle: 0,
};

const messageBoxDefaults: Partial<Settings['messageBoxSettings']> = {
	enabled: false,
	padding: 0,
	borderRadius: 0,
};

const spitterEmoteDefaults: Partial<Settings['spitterEmoteSettings']> = {
	enabled: false,
};

export const useDudesSocket = defineStore('dudes-socket', () => {
	const dudesStore = useDudes();
	const { dudes } = storeToRefs(dudesStore);
	// сюда добавил просто для демонстрации как с этим работать. По сути нам нужно каждый раз очищать этот массив
	// когда мы коннектимся к сокету, потому что сокет каждый раз отправляет все настройки всех юзеров при коннекте
	// смотри onConnected
	const usersSettings = ref<DudesUserSettings[]>([]);

	const { updateSettings, updateChannelData, loadFont } = useDudesSettings();
	const overlayId = ref('');
	const dudesUrl = ref('');
	const { data, send, open, close, status } = useWebSocket(
		dudesUrl,
		{
			immediate: false,
			autoReconnect: {
				delay: 500,
			},
			onConnected() {
				send(JSON.stringify({ eventName: 'getSettings' }));
				usersSettings.value = [];
			},
		},
	);

	watch(data, async (recieviedData) => {
		if (!dudes.value) return;

		const parsedData = JSON.parse(recieviedData) as TwirWebSocketEvent;
		if (parsedData.eventName === 'settings') {
			const data = parsedData.data as Required<Settings & ChannelData>;

			updateChannelData({
				channelId: data.channelId,
				channelName: data.channelName,
				channelDisplayName: data.channelDisplayName,
			});

			updateSettingFromSocket(data);
		}

		if (parsedData.eventName === 'jump') {
			const data = parsedData.data as DudesJumpRequest;
			const dude = dudes.value.getDude(data.userDisplayName);
			if (dude) {
				dudesStore.jumpDude(data);
			} else {
				dudesStore.createDude(data.userDisplayName, data.userColor);
			}
		}

		if (parsedData.eventName === 'punished') {
			const data = parsedData.data as DudesUserPunishedRequest;
			dudes.value.removeDude(data.userDisplayName);
		}

		// for each 100 users, we get a new event with 100 settings entities
		if (parsedData.eventName === 'usersSettings') {
			const data = parsedData.data as DudesUserSettings[];
			usersSettings.value = [...usersSettings.value, ...data];
		}

		if (parsedData.eventName === 'changeColor') {
			const data = parsedData.data as DudesChangeColorRequest;

			const dude = dudes.value.getDude(data.userName);
			if (dude) {
				dude.bodyTint(data.color);
			} else {
				dudesStore.createDude(data.userName, data.color);
			}
		}

		if (parsedData.eventName === 'grow') {
			const data = parsedData.data as DudesGrowRequest;
			const dude = dudes.value.getDude(data.userName);
			if (dude) {
				dude.grow();
			} else {
				dudesStore.createDude(data.userName, data.color);
			}
		}
	});

	async function updateSettingFromSocket(data: Required<Settings>) {
		const fontFamily = await loadFont(
			data.nameBoxSettings.fontFamily,
			data.nameBoxSettings.fontWeight,
			data.nameBoxSettings.fontStyle,
		);

		updateSettings({
			ignore: data.ignoreSettings,
			dudes: {
				dude: {
					...data.dudeSettings,
					sounds: {
						...soundsDefaults,
						enabled: data.dudeSettings.soundsEnabled,
						volume: data.dudeSettings.soundsVolume,
					},
				},
				name: {
					...nameBoxDefaults,
					...data.nameBoxSettings,
					fontFamily,
				},
				message: {
					...messageBoxDefaults,
					...data.messageBoxSettings,
					fontFamily,
				},
				spitter: {
					...spitterEmoteDefaults,
					...data.spitterEmoteSettings,
				},
			},
		});
	}

	function destroy(): void {
		if (status.value === 'OPEN') {
			close();
		}
	}

	function connect(apiKey: string, id: string): void {
		overlayId.value = id;
		dudesUrl.value = generateSocketUrlWithParams('/overlays/dudes', {
			apiKey,
			id,
		});
		open();
	}

	return {
		destroy,
		connect,
		updateSettingFromSocket,
	};
});
