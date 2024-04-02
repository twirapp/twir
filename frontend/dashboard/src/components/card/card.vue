<script setup lang="ts">
import { NCard, useThemeVars, NSpace } from 'naive-ui';
import { FunctionalComponent, computed } from 'vue';

const themeVars = useThemeVars();
const titleColor = computed(() => themeVars.value.textColor1);

withDefaults(defineProps<{
	title: string,
	icon?: FunctionalComponent,
	iconStroke?: number,
	withStroke?: boolean,
	iconFill?: string
	iconWidth?: string
	iconHeight?: string
	isLoading?: boolean
}>(), {
	withStroke: true,
	iconWidth: '48px',
	iconHeight: '48px',
});

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
		<div class="flex flex-col flex-1 h-full">
			<component
				:is="icon"
				v-if="icon"
				:style="{
					color: iconFill,
					fill: iconFill ? 'currentColor' : null,
					stroke: withStroke ? '#61e8bb' : null,
					strokeWidth: iconStroke,
					width: iconWidth,
					height: iconHeight,
					marginBottom: '16px'
				}"
			/>
			<n-space justify="space-between">
				<h2 class="text-xl mb-3" :style="{ color: titleColor }">
					{{ title }}
				</h2>
				<slot name="headerExtra" />
			</n-space>
			<div :style="{ color: themeVars.textColor3, 'margin-bottom': '10px' }">
				<slot name="content" />
			</div>
			<div class="footer flex gap-2 mt-auto flex-wrap">
				<slot name="footer" />
			</div>
		</div>
	</n-card>
</template>

<style scoped>
.footer :deep(button span) {
	@apply text-sm;
}

.footer :deep(button svg) {
	@apply h-5 w-5;
}

@media (max-width: 568px) {
	.footer :deep(button) {
		@apply w-full;
	}
}
</style>
