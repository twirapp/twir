<script setup lang="ts">
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import { DraggableResizableContainer, DraggableResizableVue } from 'draggable-resizable-vue3';

import Bot from '@/components/dashboard/bot.vue';
import Chat from '@/components/dashboard/chat.vue';
import { usePositions } from '@/components/dashboard/positions';
import Stats from '@/components/dashboard/stats.vue';
import Stream from '@/components/dashboard/stream.vue';

const widgetsPositions = usePositions();
</script>

<template>
	<Stats />
	<draggable-resizable-container
		:grid="[40, 40]"
		:show-grid="true"
		class="drag-container"
	>
		<draggable-resizable-vue
			v-model:x="widgetsPositions.chat.x"
			v-model:y="widgetsPositions.chat.y"
			v-model:h="widgetsPositions.chat.height"
			v-model:w="widgetsPositions.chat.width"
			v-model:active="widgetsPositions.chat.isActive"
			:minWidth="300"
			:minHeight="400"
		>
			<Chat />
		</draggable-resizable-vue>

		<draggable-resizable-vue
			v-model:x="widgetsPositions.botManage.x"
			v-model:y="widgetsPositions.botManage.y"
			v-model:h="widgetsPositions.botManage.height"
			v-model:w="widgetsPositions.botManage.width"
			v-model:active="widgetsPositions.botManage.isActive"
			:minWidth="330"
			:minHeight="200"
		>
			<Bot />
		</draggable-resizable-vue>

		<draggable-resizable-vue
			v-model:x="widgetsPositions.stream.x"
			v-model:y="widgetsPositions.stream.y"
			v-model:h="widgetsPositions.stream.height"
			v-model:w="widgetsPositions.stream.width"
			v-model:active="widgetsPositions.stream.isActive"
			:minWidth="330"
			:minHeight="200"
		>
			<Stream />
		</draggable-resizable-vue>
	</draggable-resizable-container>
</template>

<style scoped>
.drag-container {
	-webkit-user-select: none;
	-ms-user-select: none;
	user-select: none;
	width: auto;
	padding: 10px;
	height: 79vh;
}

</style>
