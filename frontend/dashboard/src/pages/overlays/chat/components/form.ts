import type {
	Settings,
} from '@twir/grpc/generated/api/api/overlays_chat';
import { ref, toRaw } from 'vue';

type SettingsWithOptionalId = Omit<Settings, 'id'> & { id?: string }

const defaultSettings: SettingsWithOptionalId = {
	fontFamily: 'inter',
	fontSize: 20,
	fontWeight: 400,
	fontStyle: 'normal',
	hideBots: false,
	hideCommands: false,
	messageHideTimeout: 0,
	messageShowDelay: 0,
	preset: 'clean',
	showBadges: true,
	showAnnounceBadge: true,
	textShadowColor: 'rgba(0,0,0,1)',
	textShadowSize: 0,
	chatBackgroundColor: 'rgba(0, 0, 0, 0)',
	direction: 'top',
};

const data = ref<SettingsWithOptionalId>(structuredClone(defaultSettings));

export const useChatOverlayForm = () => {
	function $setData(d: SettingsWithOptionalId) {
		data.value = structuredClone(toRaw(d));
	}

	function $reset() {
		data.value = structuredClone(defaultSettings);
	}

	function $getDefaultSettings() {
		return structuredClone(defaultSettings);
	}

	return {
		data,
		$setData,
		$reset,
		$getDefaultSettings,
	};
};
