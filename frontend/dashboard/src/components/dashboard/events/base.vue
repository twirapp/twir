<script lang="ts" setup>
import type { SVGProps } from '@tabler/icons-vue';
import { UseTimeAgo } from '@vueuse/components';
import type { FunctionalComponent } from 'vue';

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
	<div class="flex min-h-[50px] gap-2.5 px-1 select-text border-b-[color:var(--n-border-color)] border-b border-solid">
		<div class="flex justify-between items-center w-full">
			<div class="flex gap-2.5 items-center">
				<component :is="icon" class="flex items-center h-9 w-9" :style="{ color: iconColor }" />
				<div class="flex flex-col">
					<slot name="leftContent" />
				</div>
			</div>

			<div class="flex items-end text-xs h-full px-2.5 flex-shrink-0">
				<UseTimeAgo v-slot="{ timeAgo }" :time="new Date(Number(createdAt))" :update-interval="1000" show-second>
					{{ timeAgo }}
				</UseTimeAgo>
			</div>
		</div>
	</div>
</template>
