<script setup lang='ts'>
import { IconCalendarPlus } from '@tabler/icons-vue';
import { NCard, NGrid, NGridItem, NInput, NModal, NSkeleton, useThemeVars } from 'naive-ui';
import { computed, ref, toRaw } from 'vue';

import { useEventsManager, useUserAccessFlagChecker } from '@/api/index.js';
import Card from '@/components/events/card.vue';
import { getEventName } from '@/components/events/helpers.js';
import Modal from '@/components/events/modal.vue';
import type { EditableEvent } from '@/components/events/types.js';
import { ChannelRolePermissionEnum } from '@/gql/graphql';

const themeVars = useThemeVars();
const cardHoverColor = computed(() => themeVars.value.hoverColor);

const eventsManager = useEventsManager();
const { data: eventsList, isLoading } = eventsManager.getAll({});

const showModal = ref(false);
const editableEvent = ref<EditableEvent | null>(null);

const userCanManageEvents = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageEvents);

function openSettings(id?: string) {
	if (!userCanManageEvents.value) return;

	const event = eventsList.value?.events.find(e => e.id === id);
	if (event) {
		editableEvent.value = structuredClone(toRaw(event));
	} else {
		editableEvent.value = null;
	}

	showModal.value = true;
}

const search = ref('');
const events = computed(() => {
	const s = search.value.toLocaleLowerCase();

	return eventsList.value?.events.filter((e) => {
		return getEventName(e.type).toLocaleLowerCase().includes(s)
			|| e.description.toLocaleLowerCase().includes(s);
	}) ?? [];
});
</script>

<template>
	<div>
		<n-grid v-if="isLoading" responsive="screen" cols="1 s:2 m:2 l:3" :x-gap="10" :y-gap="10">
			<n-grid-item v-for="index in 6" :key="index" :span="1">
				<n-card content-style="padding: 0px">
					<n-skeleton height="300px" />
				</n-card>
			</n-grid-item>
		</n-grid>

		<n-input v-model:value="search" clearable placeholder="Search..." class="w-[30%]" />

		<n-grid
			v-if="!isLoading"
			style="margin-top: 15px;"
			responsive="screen"
			item-responsive cols="1 s:2 m:2 l:3"
			:x-gap="10"
			:y-gap="10"
		>
			<n-grid-item>
				<n-card
					class="new-event-card h-full flex items-center justify-center"
					:style="{ cursor: userCanManageEvents ? 'pointer' : 'not-allowed' }"
					embedded
					@click="openSettings"
				>
					<IconCalendarPlus class="w-20 h-20" />
				</n-card>
			</n-grid-item>
			<n-grid-item v-for="event of events" :key="event.id">
				<card :event="event" @open-settings="openSettings(event.id)" />
			</n-grid-item>
		</n-grid>
	</div>

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
.new-event-card:hover {
	background-color: v-bind(cardHoverColor);
}
</style>
