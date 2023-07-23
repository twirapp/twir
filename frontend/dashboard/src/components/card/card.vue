<script setup lang="ts">
import { NCard, useThemeVars, NSpace } from 'naive-ui';
import { FunctionalComponent, computed } from 'vue';

const themeVars = useThemeVars();
const titleColor = computed(() => themeVars.value.textColor1);

defineProps<{
	title: string;
	icon?: FunctionalComponent;
}>();

defineEmits<{
	openSettings: [];
}>();

defineSlots<{
	content?: FunctionalComponent
	footer?: FunctionalComponent
	headerExtra?: FunctionalComponent
}>();

</script>

<template>
	<n-card embedded>
		<div style="display:flex; flex-direction:column; flex:1; height: 100%">
			<component
				:is="icon"
				v-if="icon"
				style="width: 48px; height: 48px; stroke-width: 2px; stroke: #61e8bb; margin-bottom: 16px"
			/>
			<n-space justify="space-between">
				<h2 class="card-title">
					{{ title }}
				</h2>
				<slot name="headerExtra" />
			</n-space>
			<div :style="{ color: themeVars.textColor3, 'margin-bottom': '10px' }">
				<slot name="content" />
			</div>
			<div class="footer" style="margin-top: auto;">
				<slot name="footer" />
			</div>
		</div>
	</n-card>
</template>

<style scoped>
.card-title {
	color: v-bind(titleColor);
	margin: 0 0 12px 0;
	font-size: 20px;
	line-height: 24px;
}

.footer {
	display: flex;
	gap: 8px;
	margin-top: 20px;
	flex-wrap: wrap;
}

.footer :deep(button span) {
	font-size: 14px;
	line-height: 20px
}

.footer :deep(button svg) {
	height: 20px;
	width: 20px;
	margin-left: 8px
}

@media (max-width: 568px) {
	.footer :deep(button) {
		width: 100%;
	}
}
</style>
