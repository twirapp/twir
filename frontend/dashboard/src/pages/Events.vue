<script setup lang='ts'>
import { IconCalendarPlus } from '@tabler/icons-vue';
import { NCard, NGrid, NGridItem, NSkeleton, useThemeVars, NModal } from 'naive-ui';
import { computed, ref, toRaw } from 'vue';

import { useEventsManager } from '@/api/index.js';
import Card from '@/components/events/card.vue';
import Modal from '@/components/events/modal.vue';
import type { EditableEvent } from '@/components/events/types.js';

const themeVars = useThemeVars();
const cardHoverColor = computed(() => themeVars.value.hoverColor);

const eventsManager = useEventsManager();
const { data: eventsList, isLoading } = eventsManager.getAll({});

const showModal = ref(false);
const editableEvent = ref<EditableEvent | null>(null);

function openSettings(id?: string) {
	const event = eventsList.value?.events.find(e => e.id === id);
	if (event) {
		editableEvent.value = structuredClone(toRaw(event));
	} else {
		editableEvent.value = null;
	}

	showModal.value = true;
}
</script>

<template>
	<!-- <Transition appear mode="out-in"> -->
	<n-grid v-if="isLoading" responsive="screen" cols="1 s:2 m:2 l:3" :x-gap="10" :y-gap="10">
		<n-grid-item v-for="index in 6" :key="index" :span="1">
			<n-card content-style="padding: 0px">
				<n-skeleton height="300px" />
			</n-card>
		</n-grid-item>
	</n-grid>

	<n-grid v-else responsive="screen" item-responsive cols="1 s:2 m:2 l:3" :x-gap="10" :y-gap="10">
		<n-grid-item>
			<n-card
				class="new-event-card"
				content-style="
						display: flex;
						align-items: center;
						justify-content: center;
					"
				embedded
				@click="openSettings"
			>
				<IconCalendarPlus style="height: 80px; width: 80px;" />
			</n-card>
		</n-grid-item>
		<n-grid-item v-for="event of eventsList!.events" :key="event.id">
			<card :event="event" @open-settings="openSettings(event.id)" />
		</n-grid-item>
	</n-grid>
	<!-- </Transition> -->

	<n-modal
		v-model:show="showModal"
		:mask-closable="false"
		:segmented="true"
		preset="card"
		:title="editableEvent?.type ?? 'New event'"
		class="modal"
		:style="{
			width: '800px',
		}"
		@close="editableEvent = null"
	>
		<modal :event="editableEvent" @saved="showModal = false" />
	</n-modal>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: all 0.5s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
  transform: scale(0.9);
}

.new-event-card {
	/* height: 120px; */
	cursor: pointer;
	height: 100%;
}

.new-event-card:hover {
	background-color: v-bind(cardHoverColor);
}
</style>
