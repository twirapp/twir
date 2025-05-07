<script setup lang='ts'>
import { IconCalendarPlus } from '@tabler/icons-vue'
import { NCard, NGrid, NGridItem, NInput, NModal, NSkeleton, useThemeVars } from 'naive-ui'
import { computed, ref, toRaw } from 'vue'

import type { EditableEvent } from '@/features/events/constants/types.js'

import { useEventsManager, useUserAccessFlagChecker } from '@/api/index.js'
import Card from '@/features/events/constants/card.vue'
import { getEventName } from '@/features/events/constants/helpers.js'
import Modal from '@/features/events/constants/modal.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const themeVars = useThemeVars()
const cardHoverColor = computed(() => themeVars.value.hoverColor)

const eventsManager = useEventsManager()
const { data: eventsList, isLoading } = eventsManager.getAll({})

const showModal = ref(false)
const editableEvent = ref<EditableEvent | null>(null)

const userCanManageEvents = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageEvents)

function openSettings(id?: string) {
	if (!userCanManageEvents.value) return

	const event = eventsList.value?.events.find(e => e.id === id)
	if (event) {
		editableEvent.value = structuredClone(toRaw(event))
	} else {
		editableEvent.value = null
	}

	showModal.value = true
}

const search = ref('')
const events = computed(() => {
	const s = search.value.toLocaleLowerCase()

	return eventsList.value?.events.filter((e) => {
		return getEventName(e.type).toLocaleLowerCase().includes(s)
		  || e.description.toLocaleLowerCase().includes(s)
	}) ?? []
})
</script>

<template>
	<div>
		<NGrid v-if="isLoading" responsive="screen" cols="1 s:2 m:2 l:3" :x-gap="10" :y-gap="10">
			<NGridItem v-for="index in 6" :key="index" :span="1">
				<NCard content-style="padding: 0px">
					<NSkeleton height="300px" />
				</NCard>
			</NGridItem>
		</NGrid>

		<NInput v-model:value="search" clearable placeholder="Search..." class="w-[30%]" />

		<NGrid
			v-if="!isLoading"
			style="margin-top: 15px;"
			responsive="screen"
			item-responsive cols="1 s:2 m:2 l:3"
			:x-gap="10"
			:y-gap="10"
		>
			<NGridItem>
				<NCard
					class="new-event-card h-full flex items-center justify-center"
					:style="{ cursor: userCanManageEvents ? 'pointer' : 'not-allowed' }"
					embedded
					@click="openSettings"
				>
					<IconCalendarPlus class="w-20 h-20" />
				</NCard>
			</NGridItem>
			<NGridItem v-for="event of events" :key="event.id">
				<Card :event="event" @open-settings="openSettings(event.id)" />
			</NGridItem>
		</NGrid>
	</div>

	<NModal
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
		<Modal :event="editableEvent" @saved="showModal = false" />
	</NModal>
</template>

<style scoped>
.new-event-card:hover {
	background-color: v-bind(cardHoverColor);
}
</style>
