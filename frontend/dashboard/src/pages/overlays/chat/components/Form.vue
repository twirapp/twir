<script setup lang="ts">
import { IconReload } from '@tabler/icons-vue';
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
import { ref, computed, watch } from 'vue';
import { useI18n } from 'vue-i18n';

import { defaultChatSettings } from './default-settings';
import { useChatOverlayForm } from './form.js';

import { useChatOverlayManager, useProfile, useUserAccessFlagChecker } from '@/api';
import { useCopyOverlayLink } from '@/components/overlays/copyOverlayLink.js';
import { useNaiveDiscrete } from '@/composables/use-naive-discrete';
import { storeToRefs } from 'pinia';

const { t } = useI18n();
const themeVars = useThemeVars();
const discrete = useNaiveDiscrete();
const { copyOverlayLink } = useCopyOverlayLink('chat');
const userCanEditOverlays = useUserAccessFlagChecker('MANAGE_OVERLAYS');
const { data: profile } = storeToRefs(useProfile());

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
watch(() => fontData.value, (font) => {
	if (!font) return;
	formValue.value.fontFamily = font.id;
});

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
		<div class="flex-wrap justify-start">
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
					<n-text class="text-xd mt-1">
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

				<div class="slider">
					<span>{{ t('overlays.chat.paddingContainer') }} ({{ formValue.paddingContainer }}px)</span>
					<n-slider
						v-model:value="formValue.paddingContainer" :min="0" :max="256"
						:marks="{ 0: '0', 256: '256'}"
					/>
				</div>

				<n-divider />
				<div>
					<span>{{ t('overlays.chat.fontFamily') }}</span>
					<font-selector
						v-model:font="fontData"
						:font-family="formValue.fontFamily"
						:font-weight="formValue.fontWeight"
						:font-style="formValue.fontStyle"
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
					<div class="flex justify-between mb-1">
						<span>{{ t('overlays.chat.backgroundColor') }}</span>
						<n-button
							size="tiny" secondary type="success"
							@click="formValue.chatBackgroundColor = defaultChatSettings.chatBackgroundColor"
						>
							<IconReload class="h-4 w-4" />
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

.card-body-column {
	width: 100%;
}

.switch {
	@apply flex justify-between;
}

.card {
	background-color: v-bind('themeVars.cardColor');
}
</style>
