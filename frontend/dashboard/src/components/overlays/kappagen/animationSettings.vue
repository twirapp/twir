<script lang="ts" setup>
import { IconPlayerPlay } from '@tabler/icons-vue';
import type { Settings_AnimationSettings } from '@twir/grpc/generated/api/api/overlays_kappagen';
import { NSwitch, NButton, NInputNumber, useThemeVars, NDynamicInput } from 'naive-ui';
import { computed } from 'vue';

import { animations } from './kappagen_animations';

// eslint-disable-next-line no-undef
const settings = defineModel<Settings_AnimationSettings>('settings');

const animation = computed(() => {
	return animations.find(a => a.style === settings.value?.style);
});

defineEmits<{
	play: [animation: Settings_AnimationSettings]
}>();

const themeVars = useThemeVars();
</script>

<template>
	<div v-if="animation" class="card">
		<div class="content">
			<div
				class="title"
				:class="{
					'title-bordered': animation.count !== undefined || animation.prefs != undefined
				}"
			>
				<div class="info">
					<span>{{ animation.style }}</span>
					<n-button text size="tiny" @click="$emit('play', settings)">
						<IconPlayerPlay style="display: flex; height: 18px;" />
					</n-button>
				</div>
				<n-switch v-model:value="settings!.enabled" />
			</div>

			<div class="settings">
				<div v-if="animation.prefs?.size !== undefined" class="form-item">
					<span>Size</span>
					<n-input-number
						v-model:value="settings!.prefs!.size"
						:step="animation.style === 'TheCube' ? 0.1 : 1"
						size="tiny"
					/>
				</div>

				<div v-if="animation.prefs?.center !== undefined" class="form-item">
					<span>Center</span>
					<n-switch v-model:value="settings!.prefs!.center" />
				</div>

				<div v-if="animation.prefs?.faces !== undefined" class="form-item">
					<span>Faces</span>
					<n-switch v-model:value="settings!.prefs!.faces" />
				</div>

				<div v-if="animation.prefs?.speed !== undefined" class="form-item">
					<span>Speed</span>
					<n-input-number
						v-model:value="settings!.prefs!.speed"
						size="tiny"
						:max="100"
					/>
				</div>


				<div v-if="animation.count !== undefined" class="form-item">
					<span>Emotes count</span>
					<n-input-number
						v-model:value="settings!.count"
						:min="1"
						:max="150"
						size="tiny"
					/>
				</div>

				<div v-if="animation.prefs?.message !== undefined" style="display: flex; flex-direction: column; gap: 4px;">
					<span>Texts</span>
					<n-dynamic-input
						v-model:value="settings!.prefs!.message"
						:max="10"
					/>
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.card {
	display: flex;
	flex-direction: column;
	gap: 8px;
	height: 100%;
	border-radius: 4px;
	background-color: v-bind('themeVars.actionColor');
}

.card .content {
	padding: 6px;
}

.card .content .settings {
	padding-top: 5px;
	display: flex;
	flex-direction: column;
	gap: 8px;
}

.card .title {
	display: flex;
	justify-content: space-between;
	width: 100%;
	padding-bottom: 3px;
}

.card .title-bordered {
	border-bottom: 1px solid v-bind('themeVars.borderColor');
}

.card .title .info {
	display: flex;
	gap: 4px;
}

.card .form-item {
	display: flex;
	justify-content: space-between;
	gap: 4px;
}

:deep(.n-input-number) {
	width: 20%
}
</style>
