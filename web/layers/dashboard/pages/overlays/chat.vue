<script setup lang="ts">
import { until } from '@vueuse/core'
import { useRouteQuery } from '@vueuse/router'
import { type Component, computed, ref, watch } from 'vue'
import { useChatOverlayForm } from '~~/layers/dashboard/pages/overlays/chat/components/form.js'
import Form from '~~/layers/dashboard/pages/overlays/chat/components/Form.vue'

import { useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth'
import { useChatOverlayApi } from '~~/layers/dashboard/api/overlays/chat.js'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'
import { Button } from '@/components/ui/button'
import {
	AlertDialog,
	AlertDialogAction,
	AlertDialogCancel,
	AlertDialogContent,
	AlertDialogDescription,
	AlertDialogFooter,
	AlertDialogHeader,
	AlertDialogTitle,
} from '@/components/ui/alert-dialog'
import PageLayout, { type PageLayoutTab } from '~~/layers/dashboard/layout/page-layout.vue'

definePageMeta({ layout: 'dashboard', middleware: 'auth' })

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

		if (!activeTab.value) {
			activeTabQuery.value = newValue[0].id!
			return
		}

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

	const remainingCount = overlays.length - 1
	if (remainingCount === 0) {
		activeTabQuery.value = null
	} else if (currentIndex === 0) {
		activeTabQuery.value = overlays[1]?.id ?? null
	} else {
		activeTabQuery.value = overlays[0]?.id ?? null
	}

	deleteDialogOpen.value = false
	presetToDelete.value = undefined
}

async function handleAdd() {
	const previousLength = chatOverlaysData.value?.chatOverlays.length ?? 0

	await creator.executeMutation({ input: getDefaultSettings() })

	await until(() => chatOverlaysData.value?.chatOverlays).changed()

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
				<Button
					v-if="canAddMore"
					size="sm"
					variant="default"
					@click="handleAdd"
				>
					<Icon name="lucide:plus" class="h-4 w-4 mr-2" />
					{{ t('sharedButtons.add') || 'Add Preset' }}
				</Button>
				<Button
					v-if="currentTab && userCanEditOverlays"
					size="sm"
					variant="destructive"
					@click="handleDeleteClick(currentTab)"
				>
					<Icon name="lucide:trash2" class="h-4 w-4 mr-2" />
					{{ t('sharedButtons.delete') || 'Delete' }}
				</Button>
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
					<Button
						v-if="userCanEditOverlays"
						size="lg"
						variant="default"
						@click="handleAdd"
					>
						<Icon name="lucide:plus" class="h-5 w-5 mr-2" />
						{{ t('sharedButtons.create') || 'Create Preset' }}
					</Button>
				</div>
			</div>
		</template>
	</PageLayout>

	<AlertDialog v-model:open="deleteDialogOpen">
		<AlertDialogContent>
			<AlertDialogHeader>
				<AlertDialogTitle>{{ t('overlays.chat.deleteTitle') || 'Delete Preset' }}</AlertDialogTitle>
				<AlertDialogDescription>
					{{ t('overlays.chat.deleteDescription') || 'Are you sure you want to delete this preset? This action cannot be undone.' }}
				</AlertDialogDescription>
			</AlertDialogHeader>
			<AlertDialogFooter>
				<AlertDialogCancel>{{ t('sharedButtons.cancel') || 'Cancel' }}</AlertDialogCancel>
				<AlertDialogAction @click="confirmDelete">
					{{ t('sharedButtons.delete') || 'Delete' }}
				</AlertDialogAction>
			</AlertDialogFooter>
		</AlertDialogContent>
	</AlertDialog>
</template>
