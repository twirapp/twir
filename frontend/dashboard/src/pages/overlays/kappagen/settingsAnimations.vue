<script setup lang="ts">
import { IconPlayerPlay } from '@tabler/icons-vue';
import type { Settings_AnimationSettings } from '@twir/api/messages/overlays_kappagen/overlays_kappagen';
import { NButton, NDynamicInput, NGrid, NGridItem, NInputNumber, NSwitch } from 'naive-ui';
import { watch } from 'vue';

import { animations } from './kappagen_animations';
import { useKappagenFormSettings } from './store.js';

const { settings: formValue } = useKappagenFormSettings();

defineEmits<{
	play: [animation: Settings_AnimationSettings]
}>();

watch(formValue.value.animations, (v) => {
	for (const animation of animations) {
		const exists = v.find(a => a.style === animation.style);
		if (exists) continue;

		formValue.value.animations.push(animation);
	}
}, { immediate: true });
</script>

<template>
	<n-grid :cols="2" :x-gap="16" :y-gap="16" responsive="self">
		<n-grid-item v-for="animation of formValue.animations" :key="animation.style" :span="1">
			<div class="card">
				<div class="content">
					<div
						class="title"
						:class="{
							'title-bordered': animation.count !== undefined || animation.prefs !== undefined
						}"
					>
						<div class="info">
							<span>{{ animation.style }}</span>
							<n-button text size="tiny" @click="$emit('play', animation)">
								<IconPlayerPlay class="flex h-4" />
							</n-button>
						</div>
						<n-switch v-model:value="animation.enabled" />
					</div>

					<div class="settings">
						<div v-if="animation.prefs?.size !== undefined" class="form-item">
							<span>Size</span>
							<n-input-number
								v-model:value="animation.prefs!.size"
								:step="animation.style === 'TheCube' ? 0.1 : 1"
								size="tiny"
							/>
						</div>

						<div v-if="animation.prefs?.center !== undefined" class="form-item">
							<span>Center</span>
							<n-switch v-model:value="animation.prefs!.center" />
						</div>

						<div v-if="animation.prefs?.faces !== undefined" class="form-item">
							<span>Faces</span>
							<n-switch v-model:value="animation.prefs!.faces" />
						</div>

						<div v-if="animation.prefs?.speed !== undefined" class="form-item">
							<span>Speed</span>
							<n-input-number
								v-model:value="animation.prefs!.speed"
								size="tiny"
								:max="100"
							/>
						</div>


						<div v-if="animation.count !== undefined" class="form-item">
							<span>Emotes count</span>
							<n-input-number
								v-model:value="animation.count"
								:min="1"
								:max="150"
								size="tiny"
							/>
						</div>

						<div v-if="animation.style === 'Text'" class="flex flex-col gap-1">
							<span>Texts</span>
							<n-dynamic-input
								v-model:value="animation.prefs!.message"
								:max="10"
							/>
						</div>
					</div>
				</div>
			</div>
		</n-grid-item>
	</n-grid>
</template>
