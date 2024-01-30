import type { DudesSettings } from '@twirapp/dudes/types';
import { defineStore } from 'pinia';
import { reactive, ref } from 'vue';

export const useDudesSettings = defineStore('dudes-settings', () => {
	const dudesSettings = reactive<DudesSettings>({
		dude: {
			color: '#969696',
			maxLifeTime: 1000 * 60 * 30,
			gravity: 400,
			scale: 4,
			sounds: {
				enabled: true,
				volume: 0.01,
			},
		},
		messageBox: {
			borderRadius: 10,
			boxColor: '#eeeeee',
			fontFamily: 'roboto',
			fontSize: 20,
			padding: 10,
			showTime: 5 * 1000,
			fill: '#333333',
		},
		nameBox: {
			fontFamily: 'roboto',
			fontSize: 18,
			fill: '#ffffff',
			lineJoin: 'round',
			strokeThickness: 4,
			stroke: '#000000',
			fillGradientStops: [0],
			fillGradientType: 0,
			fontStyle: 'normal',
			fontVariant: 'normal',
			fontWeight: 400,
			dropShadow: false,
			dropShadowAlpha: 1,
			dropShadowAngle: 0,
			dropShadowBlur: 1,
			dropShadowDistance: 1,
			dropShadowColor: '#3ac7d9',
		},
	});

	const channelInfo = ref<{ channelId: string, channelName: string }>();

	function updateSettings(settings: Record<string, any>) {
		if (settings.channelId && settings.channelName) {
			channelInfo.value = {
				channelId: settings.channelId,
				channelName: settings.channelName,
			};
		}

		dudesSettings.dude = {
			color: settings.dudeColor,
			maxLifeTime: settings.dudeMaxLifeTime,
			gravity: settings.dudeGravity,
			scale: settings.dudeScale,
			sounds: {
				enabled: settings.dudeSoundsEnabled,
				volume: settings.dudeSoundsVolume,
			},
		};

		dudesSettings.messageBox = {
			borderRadius: settings.messageBoxBorderRadius,
			boxColor: settings.messageBoxBoxColor,
			fontFamily: settings.nameBoxFontSize, // TODO: change target
			fontSize: settings.messageBoxFontSize,
			padding: settings.messageBoxPadding,
			showTime: settings.messageBoxShowTime,
			fill: settings.messageBoxFill,
		};

		dudesSettings.nameBox = {
			fontFamily: settings.nameBoxFontFamily,
			fontSize: settings.nameBoxFontSize,
			fill: settings.nameBoxFill,
			lineJoin: settings.nameBoxLineJoin,
			strokeThickness: settings.nameBoxStrokeThickness,
			stroke: settings.nameBoxStroke,
			fillGradientStops: settings.nameBoxFillGradientStops,
			fillGradientType: settings.nameBoxFillGradientType,
			fontStyle: settings.nameBoxFontStyle,
			fontVariant: settings.nameBoxFontVariant,
			fontWeight: Number(settings.nameBoxFontWeight),
			dropShadow: settings.nameBoxDropShadow,
			dropShadowAlpha: settings.nameBoxDropShadowAlpha,
			dropShadowAngle: settings.nameBoxDropShadowAngle,
			dropShadowBlur: settings.nameBoxDropShadowBlur,
			dropShadowDistance: settings.nameBoxDropShadowDistance,
			dropShadowColor: settings.nameBoxDropShadowColor,
		};
	}

	return {
		channelInfo,
		dudesSettings,
		updateSettings,
	};
});
