<script lang="ts" setup>
import { type SVGProps } from '@tabler/icons-vue';
import { UseTimeAgo } from '@vueuse/components';
import { useThemeVars } from 'naive-ui';
import { computed, type FunctionalComponent } from 'vue';

const theme = useThemeVars();
const borderColor = computed(() => theme.value.borderColor);

defineProps<{
	icon: FunctionalComponent<SVGProps, Record<string, any>, any>
	iconColor?: string;
	createdAt: string
}>();

defineSlots<{
	leftContent: any,
	rightContent: any,
}>();
</script>

<template>
	<div class="event">
		<div class="content">
			<div style="display: flex; gap: 10px; align-items: center;">
				<component :is="icon" class="icon" :style="{ color: iconColor }" />
				<div style="display: flex; flex-direction: column;">
					<slot name="leftContent" />
				</div>
			</div>

			<div style="display: flex; align-items: flex-end; font-size: 11px; height: 100%; padding-right: 10px">
				<UseTimeAgo v-slot="{ timeAgo }" :time="new Date(Number(createdAt))" :update-interval="1000" show-second>
					{{ timeAgo }}
				</UseTimeAgo>
			</div>
		</div>
	</div>
</template>

<style scoped>
.event {
	display: flex;
	min-height: 50px;
	gap: 10px;
	padding-left: 5px;
	padding-right: 5px;
	border-bottom: 1px solid v-bind(borderColor);
}

.event .content {
	display: flex;
	justify-content: space-between;
	align-items: center;
	width: 100%;
}

.icon {
	height: 35px;
	width: 35px;
	display: flex;
	align-items: center;
}
</style>
