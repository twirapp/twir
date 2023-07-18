<script setup lang='ts'>
import { IconSettings, IconCopy } from '@tabler/icons-vue';
import { NCard, NSpace, NButton, NText, useThemeVars, useMessage, NTooltip } from 'naive-ui';
import { FunctionalComponent } from 'vue';

import VectorSVG from './vector.svg?component';

const themeVars = useThemeVars();

const props = defineProps<{
	description: string
	title: string
	overlayLink?: string
}>();
defineSlots<{
	icon: FunctionalComponent<any>
}>();
defineEmits<{
	openSettings: []
}>();

const messages = useMessage();
const copyOverlayLink = () => {
	if (!props.overlayLink) return;

	navigator.clipboard.writeText(props.overlayLink);
	messages.success('Copied link url, paste it in obs as browser source');
};
</script>

<template>
  <n-card content-style="padding: 0px" class="overlay-item">
    <div style="height: 100%; display: flex">
      <div class="section-icon">
        <div class="vector">
          <VectorSVG />
        </div>

        <slot name="icon" />
      </div>

      <div class="section-info">
        <h2 style="margin:0px">
          {{ title }}
        </h2>
        <n-text :style="{ color: themeVars.textColor3, 'margin-top': '12px' }">
          {{ description }}
        </n-text>
        <n-space style="margin-top: 20px;">
          <n-button secondary size="large" @click="$emit('openSettings')">
            <n-space justify="space-between" align="center">
              <n-text>Settings</n-text>
              <IconSettings style="height: 25px" />
            </n-space>
          </n-button>
          <n-tooltip :disabled="!!overlayLink">
            <template #trigger>
              <n-button size="large" secondary type="info" :disabled="!overlayLink" @click="copyOverlayLink">
                <n-space justify="space-between" align="center">
                  <n-text>Copy overlay link</n-text>
                  <IconCopy style="height: 25px" />
                </n-space>
              </n-button>
            </template>
            You should configure overlay first
          </n-tooltip>
        </n-space>
      </div>
    </div>
  </n-card>
</template>

<style scoped>
.overlay-item {
	overflow: hidden;
	display: flex;
	height: 100%;
}

.section-icon {
	background: linear-gradient(48deg, #5557E6 0%, #3B6EEF 100%);
	display: flex;
	align-items: center;
	justify-content: center;
	overflow: hidden;
	position: relative;
	flex: 0 0 120px;
}

.section-info {
	display: flex;
	align-items: flex-start;
	flex-direction: column;
	margin: 32px
}

.vector {
	display: flex;
	align-items: center;
}

.vector svg {
	position: absolute;
	transform: rotate(167deg);
}
</style>
