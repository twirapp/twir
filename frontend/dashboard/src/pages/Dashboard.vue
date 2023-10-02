<script setup lang="ts">
import { IconPencilPlus } from '@tabler/icons-vue';
import { GridLayout, GridItem } from 'grid-layout-plus';
import { NButton, NDropdown, useThemeVars } from 'naive-ui';
import { computed } from 'vue';

import Bot from '@/components/dashboard/bot.vue';
import Chat from '@/components/dashboard/chat.vue';
import Events from '@/components/dashboard/events.vue';
import Stats from '@/components/dashboard/stats.vue';
import Stream from '@/components/dashboard/stream.vue';
import { useWidgets } from '@/components/dashboard/widgets.js';

const widgets = useWidgets();
const visibleWidgets = computed(() => widgets.value.filter((v) => v.visible));
const dropdownOptions = computed(() => {
	return widgets.value
		.filter((v) => !v.visible)
		.map((v) => ({ label: v.i, key: v.i }));
});

const addWidget = (key: string) => {
	const item = widgets.value.find(v => v.i === key);
	if (!item) return;

	const widgetsLength = visibleWidgets.value.length;

	item.visible = true;
	item.x = (widgetsLength * 2) % 12;
	item.y = widgetsLength + 12;
};

const theme = useThemeVars();
const statsBackground = computed(() => theme.value.tabColor);
</script>

<template>
	<div style="display: flex; width: 100%;" :style="{ 'background-color': statsBackground }">
		<Stats />
	</div>
	<div style="width: 100%; height: 100%; padding-left: 5px;">
		<GridLayout
			v-model:layout="widgets"
			:row-height="30"
		>
			<GridItem
				v-for="item in visibleWidgets"
				:key="item.i"
				:x="item.x"
				:y="item.y"
				:w="item.w"
				:h="item.h"
				:i="item.i"
				:min-w="item.minW"
				:min-h="item.minH"
				drag-allow-from=".widgets-draggable-handle"
			>
				<Chat v-if="item.i === 'chat'" :item="item" class="item" />
				<Bot v-if="item.i === 'bot'" :item="item" class="item" />
				<Stream v-if="item.i === 'stream'" :item="item" class="item" />
				<Events v-if="item.i === 'events'" :item="item" class="item" />
			</GridItem>
		</GridLayout>
		<div v-if="dropdownOptions.length" style="position: fixed; bottom: 10px; right: 25px">
			<n-dropdown size="huge" trigger="click" :options="dropdownOptions" @select="addWidget">
				<n-button block circle type="info" style="width: 100%; height: 100%; padding: 5px;">
					<IconPencilPlus style="width: 45px; height: 45px;" />
				</n-button>
			</n-dropdown>
		</div>
	</div>
</template>

<style scoped>
.vgl-layout {
	width: 100%
}

.item {
	height: 100%;
}
</style>
