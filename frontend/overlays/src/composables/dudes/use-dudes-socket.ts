import type { Settings } from '@twir/api/messages/overlays_dudes/overlays_dudes';
import { DudesSprite, type DudesUserSettings } from '@twir/types/overlays';
import { useWebSocket } from '@vueuse/core';
import { defineStore, storeToRefs } from 'pinia';
import { onMounted, ref, watch } from 'vue';

import { getSprite } from './dudes-config.js';
import { useDudesSettings } from './use-dudes-settings';
import { useDudes } from './use-dudes.js';

import type { TwirWebSocketEvent } from '@/api.js';
import { generateSocketUrlWithParams } from '@/helpers.js';
import type { ChannelData } from '@/types.js';

declare global {
	interface GlobalEventHandlersEventMap {
		'get-user-settings': CustomEvent<string>;
	}
}

export const useDudesSocket = defineStore('dudes-socket', () => {
	const dudesStore = useDudes();
	const { dudes } = storeToRefs(dudesStore);

	const dudesSettingsStore = useDudesSettings();
	const overlayId = ref('');
	const dudesUrl = ref('');
	const { data, send, open, close, status } = useWebSocket(dudesUrl, {
		immediate: false,
		autoReconnect: {
			delay: 500,
		},
		onConnected() {
			send(JSON.stringify({ eventName: 'getSettings' }));
		},
	});

	watch(data, async (recieviedData) => {
		if (!dudes.value?.dudes) return;

		const parsedData = JSON.parse(recieviedData) as TwirWebSocketEvent;
		if (parsedData.eventName === 'settings') {
			const data = parsedData.data as Required<Settings & ChannelData>;

			dudesSettingsStore.updateChannelData({
				channelId: data.channelId,
				channelName: data.channelName,
				channelDisplayName: data.channelDisplayName,
			});

			updateSettingFromSocket(data);
			return;
		}

		const data = parsedData.data as DudesUserSettings;
		const dude = dudes.value.dudes.getDude(data?.userId);

		if (parsedData.eventName === 'userSettings') {
			const dudeSettings = dudesSettingsStore.dudesUserSettings.get(data.userId);
			if (!dudeSettings?.userDisplayName) return;

			dudesSettingsStore.dudesUserSettings.set(data.userId, {
				...dudeSettings,
				...data,
				dudeColor: data.dudeColor ?? dudeSettings.dudeColor,
			});

			const spriteData = getSprite(
				data.dudeSprite ??
				dudesSettingsStore.dudesSettings?.overlay.defaultSprite,
			);

			const createdDude = await dudesStore.createDude({
				userId: dudeSettings.userId,
			});

			if (createdDude?.dude) {
				await createdDude.dude.updateSpriteData(spriteData);
				dudesStore.updateDudeColors(createdDude.dude, data.dudeColor);
			}
		} else if (parsedData.eventName === 'jump') {
			dude?.jump();
		} else if (parsedData.eventName === 'grow') {
			dude?.grow();
		} else if (parsedData.eventName === 'leave') {
			dude?.leave();
		} else if (parsedData.eventName === 'punished') {
			dudes.value.dudes.removeDude(data.userId);
			dudesSettingsStore.dudesUserSettings.delete(data.userId);
		}
	});

	async function updateSettingFromSocket(data: Required<Settings>) {
		const fontFamily = await dudesSettingsStore.loadFont(
			data.nameBoxSettings.fontFamily,
			data.nameBoxSettings.fontWeight,
			data.nameBoxSettings.fontStyle,
		);

		dudesSettingsStore.updateSettings({
			ignore: data.ignoreSettings,
			overlay: {
				defaultSprite: data.dudeSettings.defaultSprite as keyof typeof DudesSprite,
				maxOnScreen: data.dudeSettings.maxOnScreen,
			},
			dudes: {
				dude: {
					...data.dudeSettings,
					// TODO: rename and deprecate `eyes_color`, `cosmetics_color`
					bodyColor: data.dudeSettings.color,
				},
				sounds: {
					enabled: data.dudeSettings.soundsEnabled,
					volume: data.dudeSettings.soundsVolume,
				},
				name: {
					...data.nameBoxSettings,
					// TODO: move to nameBoxSettings
					enabled: data.dudeSettings.visibleName,
					fontFamily,
				},
				message: {
					...data.messageBoxSettings,
					fontFamily,
				},
				emote: data.spitterEmoteSettings,
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

	onMounted(() => {
		document.addEventListener('get-user-settings', (event) => {
			if (status.value !== 'OPEN') return;
			send(JSON.stringify({ eventName: 'getUserSettings', data: event.detail }));
		});
	});

	return {
		destroy,
		connect,
		updateSettingFromSocket,
	};
});
