<script setup lang="ts">
import { until } from '@vueuse/core'
import { useRouteQuery } from '@vueuse/router'
import { type Component, computed, ref, watch } from 'vue'

import { Plus, Trash2 } from 'lucide-vue-next'

import { useChatOverlayForm } from './components/form.js'
import Form from './components/Form.vue'

import { useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useChatOverlayApi } from '#layers/dashboard/api/overlays/chat'
import { ChannelRolePermissionEnum } from '~/gql/graphql'


import PageLayout, { type PageLayoutTab } from '~/layout/page-layout.vue'

const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)

const { t } = useI18n()

const chatOverlaysManager = useChatOverlayApi()
const deleter = chatOverlaysManager.useOverlayDelete()
const creator = chatOverlaysManager.useOverlayCreate()
const { data: chatOverlaysData, fetching: fetchingOverlays } = chatOverlaysManager.useOverlaysQuery()

const { setData, getDefaultSettings } = useChatOverlayForm()

const activeTabQuery = useRouteQuery<string | null>('tab', null)
const activeTab = computed<string | undefined>({
	get: () => activeTabQuery.value || undefined,
	set: (value) => {
		activeTabQuery.value = value || null
	}
})

const deleteDialogOpen = ref(false)
const presetToDelete = ref<string>()

const tabs = computed<PageLayoutTab[]>(() => {
	if (!chatOverlaysData.value?.chatOverlays.length) return []

	return chatOverlaysData.value.chatOverlays.map((overlay, index) => ({
		name: overlay.id!,
		title: `Preset #${index + 1}`,
		component: Form as Component,
	}))
})

watch(
	() => chatOverlaysData.value?.chatOverlays,
	(newValue) => {
		if (!newValue?.length) {
			activeTabQuery.value = null
			return
		}

		// Auto-open first preset if no tab is selected
		if (!activeTab.value) {
			activeTabQuery.value = newValue[0].id!
			return
		}

		// Clear invalid tabs
		if (activeTab.value && !newValue.find((o) => o.id === activeTab.value)) {
			activeTabQuery.value = null
		}
	},
	{ immediate: true }
)

watch(activeTab, (id) => {
	if (!id) return
	const entity = chatOverlaysData.value?.chatOverlays.find((s) => s.id === id)
	if (!entity) return

	setData(entity)
})

// Initialize form data when data loads (for page refresh with ?tab=id)
watch(
	() => chatOverlaysData.value?.chatOverlays,
	(overlays) => {
		if (!overlays?.length || !activeTab.value) return

		const entity = overlays.find((s) => s.id === activeTab.value)
		if (entity) {
			setData(entity)
		}
	},
	{ immediate: true }
)

function handleDeleteClick(id: string) {
	presetToDelete.value = id
	deleteDialogOpen.value = true
}

async function confirmDelete() {
	if (!presetToDelete.value) return

	const entity = chatOverlaysData.value?.chatOverlays.find((s) => s.id === presetToDelete.value)
	if (!entity?.id) return

	const overlays = chatOverlaysData.value?.chatOverlays ?? []
	const currentIndex = overlays.findIndex(o => o.id === entity.id)

	await deleter.executeMutation({ id: entity.id })

	// Navigate after deletion
	const remainingCount = overlays.length - 1
	if (remainingCount === 0) {
		// No presets left
		activeTabQuery.value = null
	} else if (currentIndex === 0) {
		// Deleted first preset, navigate to new first (was second)
		activeTabQuery.value = overlays[1]?.id ?? null
	} else {
		// Deleted other preset, navigate to first
		activeTabQuery.value = overlays[0]?.id ?? null
	}

	deleteDialogOpen.value = false
	presetToDelete.value = undefined
}

async function handleAdd() {
	const previousLength = chatOverlaysData.value?.chatOverlays.length ?? 0

	await creator.executeMutation({ input: getDefaultSettings() })

	// Wait for data to update reactively
	await until(() => chatOverlaysData.value?.chatOverlays).changed()

	// Navigate to newly created preset (last in array)
	const overlays = chatOverlaysData.value?.chatOverlays
	if (overlays && overlays.length > previousLength) {
		activeTabQuery.value = overlays[overlays.length - 1].id!
	}
}

const canAddMore = computed(() => {
	return userCanEditOverlays.value && (chatOverlaysData.value?.chatOverlays.length ?? 0) < 5
})

const hasOverlays = computed(() => {
	return (chatOverlaysData.value?.chatOverlays.length ?? 0) > 0
})


</script>

<template>
	<PageLayout :active-tab="activeTab || ''" :tabs="tabs">
		<template #title>
			{{ t('overlays.chat.title') || 'Chat Overlay' }}
		</template>

		<template #action="{ activeTab: currentTab }">
			<div class="flex gap-2">
				<UiButton
					v-if="canAddMore"
					size="sm"
					variant="default"
					@click="handleAdd"
				>
					<Plus class="h-4 w-4 mr-2" />
					{{ t('sharedButtons.add') || 'Add Preset' }}
				</UiButton>
				<UiButton
					v-if="currentTab && userCanEditOverlays"
					size="sm"
					variant="destructive"
					@click="handleDeleteClick(currentTab)"
				>
					<Trash2 class="h-4 w-4 mr-2" />
					{{ t('sharedButtons.delete') || 'Delete' }}
				</UiButton>
			</div>
		</template>

		<template v-if="fetchingOverlays" #content>
			<div class="flex flex-col items-center justify-center min-h-[60vh] gap-6">
				<div class="text-center space-y-4">
					<div class="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full mx-auto"></div>
					<p class="text-muted-foreground">Loading presets...</p>
				</div>
			</div>
		</template>

		<template v-else-if="!hasOverlays" #content>
			<div class="flex flex-col items-center justify-center min-h-[60vh] gap-6">
				<div class="text-center space-y-4 max-w-md">
					<h2 class="text-2xl font-semibold">{{ t('overlays.chat.noPresets') || 'No Presets Yet' }}</h2>
					<p class="text-muted-foreground">
						{{ t('overlays.chat.noPresetsDescription') || 'Create your first chat overlay preset to get started.' }}
					</p>
					<UiButton
						v-if="userCanEditOverlays"
						size="lg"
						variant="default"
						@click="handleAdd"
					>
						<Plus class="h-5 w-5 mr-2" />
						{{ t('sharedButtons.create') || 'Create Preset' }}
					</UiButton>
				</div>
			</div>
		</template>
	</PageLayout>

	<!-- Delete Confirmation Dialog -->
	<UiAlertDialog v-model:open="deleteDialogOpen">
		<UiAlertDialogContent>
			<UiAlertDialogHeader>
				<UiAlertDialogTitle>{{ t('overlays.chat.deleteTitle') || 'Delete Preset' }}</UiAlertDialogTitle>
				<UiAlertDialogDescription>
					{{ t('overlays.chat.deleteDescription') || 'Are you sure you want to delete this preset? This action cannot be undone.' }}
				</UiAlertDialogDescription>
			</UiAlertDialogHeader>
			<UiAlertDialogFooter>
				<UiAlertDialogCancel>{{ t('sharedButtons.cancel') || 'Cancel' }}</UiAlertDialogCancel>
				<UiAlertDialogAction @click="confirmDelete">
					{{ t('sharedButtons.delete') || 'Delete' }}
				</UiAlertDialogAction>
			</UiAlertDialogFooter>
		</UiAlertDialogContent>
	</UiAlertDialog>
</template>

<style scoped>
</style>
