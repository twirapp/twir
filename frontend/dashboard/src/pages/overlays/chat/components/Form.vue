<script setup lang="ts">
import { IconReload } from '@tabler/icons-vue';
import type {
	Settings,
} from '@twir/api/messages/overlays_chat/overlays_chat';
import { FontSelector, type Font } from '@twir/fontsource';
import {
	NButton,
	NText,
	NSwitch,
	NSlider,
	NSelect,
	NColorPicker,
	NDivider,
	useThemeVars,
} from 'naive-ui';
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';

import { useChatOverlayForm } from './form.js';

import { useChatOverlayManager, useProfile, useUserAccessFlagChecker } from '@/api';
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink.js';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete';

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

const { t } = useI18n();
const themeVars = useThemeVars();
const discrete = useNaiveDiscrete();
const { copyOverlayLink } = useCopyOverlayLink('chat');
const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');
const { data: profile } = useProfile();

const { data: formValue, $reset } = useChatOverlayForm();

const canCopyLink = computed(() => {
	return profile?.value?.selectedDashboardId === profile.value?.id && userCanEditOverlays;
});

const styleSelectOptions = [
	{ label: 'Clean', value: 'clean' },
	{ label: 'Boxed', value: 'boxed' },
];
const directionOptions = computed(() => {
	return ['top', 'right', 'bottom', 'left'].map((direction) => ({
		value: direction,
		label: t(`overlays.chat.directions.${direction}`),
	}));
});

const fontData = ref<Font | null>(null);
const fontWeightOptions = computed(() => {
	if (!fontData.value) return [];
	return fontData.value.weights.map((weight) => ({ label: `${weight}`, value: weight }));
});
const fontStyleOptions = computed(() => {
	if (!fontData.value) return [];
	return fontData.value.styles.map((style) => ({ label: style, value: style }));
});

const sliderMarks = {
	0: '0',
	60: '60',
};

const manager = useChatOverlayManager();
const updater = manager.useUpdate();

async function save() {
	if (!formValue.value.id) return;

	await updater.mutateAsync({
		id: formValue.value.id,
		settings: formValue.value,
	});

	discrete.notification.success({
		title: t('sharedTexts.saved'),
		duration: 1500,
	});
}
</script>

<template>
	<div v-if="formValue" class="card">
		<div class="card-header">
			<n-button
				secondary
				type="error"
				@click="$reset"
			>
				{{ t('sharedButtons.setDefaultSettings') }}
			</n-button>
			<n-button
				secondary
				type="info"
				:disabled="!formValue.id || !canCopyLink"
				@click="copyOverlayLink({ id: formValue.id! })"
			>
				{{ t('overlays.copyOverlayLink') }}
			</n-button>
			<n-button secondary type="success" @click="save">
				{{ t('sharedButtons.save') }}
			</n-button>
		</div>

		<div class="card-body">
			<div class="card-body-column">
				<div>
					<span>{{ t('overlays.chat.style') }}</span>
					<n-select v-model:value="formValue.preset" :options="styleSelectOptions" />
				</div>

				<div>
					<span>{{ t('overlays.chat.direction') }}</span>
					<n-select v-model:value="formValue.direction" :options="directionOptions" />
					<n-text style="font-size: 12px; margin-top: 4px;">
						{{ t('overlays.chat.directionWarning') }}
					</n-text>
				</div>

				<div class="switch">
					<span>{{ t('overlays.chat.hideBots') }}</span>
					<n-switch v-model:value="formValue.hideBots" />
				</div>

				<div class="switch">
					<span>{{ t('overlays.chat.hideCommands') }}</span>
					<n-switch v-model:value="formValue.hideCommands" />
				</div>

				<div class="switch">
					<span>{{ t('overlays.chat.showBadges') }}</span>
					<n-switch v-model:value="formValue.showBadges" />
				</div>

				<div v-if="formValue.preset === 'boxed'" class="switch">
					<span>{{ t('overlays.chat.showAnnounceBadge') }}</span>
					<n-switch
						v-model:value="formValue.showAnnounceBadge"
						:disabled="!formValue.showBadges"
					/>
				</div>

				<n-divider />
				<div>
					<span>{{ t('overlays.chat.fontFamily') }}</span>
					<font-selector
						v-model:selected-font="formValue.fontFamily"
						:font-family="formValue.fontFamily"
						:font-weight="formValue.fontWeight"
						:font-style="formValue.fontStyle"
						@update-font="(v) => fontData = v"
					/>
				</div>

				<div>
					<span>{{ t('overlays.chat.fontWeight') }}</span>
					<n-select
						v-model:value="formValue.fontWeight"
						:options="fontWeightOptions"
					/>
				</div>

				<div>
					<span>{{ t('overlays.chat.fontStyle') }}</span>
					<n-select
						v-model:value="formValue.fontStyle"
						:options="fontStyleOptions"
					/>
				</div>

				<div class="slider">
					<span>{{ t('overlays.chat.fontSize') }} ({{ formValue.fontSize }}px)</span>
					<n-slider
						v-model:value="formValue.fontSize" :min="12" :max="80"
						:marks="{ 12: '12', 80: '80'}"
					/>
				</div>

				<div class="slider">
					<div style="display: flex; justify-content: space-between; margin-bottom: 4px;">
						<span>{{ t('overlays.chat.backgroundColor') }}</span>
						<n-button
							size="tiny" secondary type="success"
							@click="formValue.chatBackgroundColor = defaultSettings.chatBackgroundColor"
						>
							<IconReload style="height: 15px;" />
							{{ t('overlays.chat.resetToDefault') }}
						</n-button>
					</div>
					<n-color-picker
						v-model:value="formValue.chatBackgroundColor"
						default-value="rgba(16, 16, 20, 1)"
					/>
				</div>


				<div class="slider">
					<span>{{ t('overlays.chat.textShadow') }}({{ formValue.textShadowSize }}px)</span>
					<n-color-picker
						v-model:value="formValue.textShadowColor"
						default-value="rgba(0,0,0,1)"
					/>
					<n-slider v-model:value="formValue.textShadowSize" :min="0" :max="30" />
				</div>

				<n-divider />

				<div class="slider">
					<span>{{ t('overlays.chat.hideTimeout') }}({{ formValue.messageHideTimeout }}s)</span>
					<n-slider
						v-model:value="formValue.messageHideTimeout" :max="60"
						:marks="sliderMarks"
					/>
				</div>

				<div class="slider">
					<span>{{ t('overlays.chat.showDelay') }}({{ formValue.messageShowDelay }}s)</span>
					<n-slider
						v-model:value="formValue.messageShowDelay" :max="60"
						:marks="sliderMarks"
					/>
				</div>
			</div>
		</div>
	</div>
</template>


<style scoped>
@import '../../styles.css';

.card-header {
	flex-wrap: wrap;
	justify-content: flex-start;
}

.card-body-column {
	width: 100%;
}

.switch {
	display: flex;
	justify-content: space-between;
}

.card {
	background-color: v-bind('themeVars.cardColor');
}
</style>
