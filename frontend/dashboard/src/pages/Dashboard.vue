<script setup lang="ts">
import { IconPencilPlus } from '@tabler/icons-vue';
import { GridLayout, GridItem } from 'grid-layout-plus';
import { NButton, NDropdown } from 'naive-ui';
import { computed, nextTick, onBeforeUnmount, onMounted, ref } from 'vue';

import BotStatus from '@/components/dashboard/bot-status.vue';
import Chat from '@/components/dashboard/chat.vue';
import Events from '@/components/dashboard/events.vue';
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

const showEmptyItem = ref(false);

const onMouseUp = () => {
	showEmptyItem.value = false;
};

onMounted(async () => {
	await nextTick();

	document.querySelectorAll('.vgl-item__resizer').forEach((el) => {
		el.addEventListener('mousedown', () => {
			showEmptyItem.value = true;
		});

		window.addEventListener('mouseup', onMouseUp);
	});
});

onBeforeUnmount(() => {
	window.removeEventListener('mouseup', onMouseUp);
});
</script>

<template>
	<bot-status />
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
				<div v-if="showEmptyItem" class="w-full h-full absolute z-50"></div>
				<Chat v-if="item.i === 'chat'" :item="item" class="item" />
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

:deep(.vgl-item__resizer) {
	z-index: 51;
}
</style>
