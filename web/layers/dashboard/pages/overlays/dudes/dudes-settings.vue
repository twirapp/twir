<script setup lang="ts">
import { computed, ref, watch } from 'vue'


import DudesSettingsForm from './dudes-settings-form.vue'
import { useDudesForm } from './use-dudes-form.js'
import { useDudesIframe } from './use-dudes-frame.js'

import type { DudesSettingsWithOptionalId } from './dudes-settings.js'

import { useProfile, useUserAccessFlagChecker } from '#layers/dashboard/api/auth'
import { useDudesOverlayManager } from '#layers/dashboard/api/overlays/dudes'




import CommandButton from '~/features/commands/ui/command-button.vue'
import { ChannelRolePermissionEnum } from '~/gql/graphql'

const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)
const dudesOverlayManager = useDudesOverlayManager()
const creator = dudesOverlayManager.useCreate()
const deleter = dudesOverlayManager.useDelete()
const { data: profile } = useProfile()

const { t } = useI18n()

const { data: entities } = dudesOverlayManager.useGetAll()

const deleteDialogOpen = ref(false)
const deleteTargetId = ref<string | null>(null)

const openedTab = ref<string>()

const { sendIframeMessage, dudesIframe } = useDudesIframe()
const dudesIframeUrl = computed(() => {
	if (!profile.value || !openedTab.value) return null
	return `${window.location.origin}/overlays/${profile.value.apiKey}/dudes?id=${openedTab.value}`
})

const { setData, getDefaultSettings } = useDudesForm()

function resetTab() {
	if (!entities.value?.dudesGetAll?.at(0)) {
		openedTab.value = undefined
		return
	}

	openedTab.value = entities.value.dudesGetAll.at(0)?.id
}

watch(
	entities,
	() => {
		resetTab()
	},
	{ immediate: true }
)

watch(openedTab, (v) => {
	const entity = entities.value?.dudesGetAll?.find((s) => s.id === v) as DudesSettingsWithOptionalId
	if (!entity) return
	setData(entity)
})

function handleClose(id: string) {
	deleteTargetId.value = id
	deleteDialogOpen.value = true
}

async function confirmDelete() {
	if (!deleteTargetId.value) return

	const entity = entities.value?.dudesGetAll?.find((s) => s.id === deleteTargetId.value)
	if (!entity?.id) return

	await deleter.executeMutation({ id: entity.id })
	resetTab()
	deleteDialogOpen.value = false
	deleteTargetId.value = null
}

async function handleAdd() {
	const { id: _id, ...rest } = getDefaultSettings()
	console.log(rest)
	await creator.executeMutation({ input: rest })
}

const addable = computed(() => {
	return userCanEditOverlays.value && (entities.value?.dudesGetAll?.length ?? 0) < 5
})
</script>

<template>
	<div class="flex gap-10" style="height: calc(100% - var(--layout-header-height))">
		<div class="relative w-[70%]">
			<div v-if="dudesIframeUrl">
				<iframe :ref="(el) => dudesIframe = el as HTMLIFrameElement" style="position: absolute" :src="dudesIframeUrl" class="iframe" />
				<div class="absolute top-4.5 left-2 flex gap-1.5">
					<UiButton size="sm" @click="sendIframeMessage('spawn-emote')"> Emote </UiButton>
					<UiButton size="sm" @click="sendIframeMessage('show-message')"> Message </UiButton>
					<UiButton size="sm" @click="sendIframeMessage('jump')"> Jump </UiButton>
					<UiButton size="sm" @click="sendIframeMessage('grow')"> Grow </UiButton>
					<UiButton size="sm" @click="sendIframeMessage('leave')"> Leave </UiButton>
					<UiButton size="sm" @click="sendIframeMessage('reset')"> Reset </UiButton>
				</div>
			</div>
		</div>
		<div class="w-[30%]">
			<div class="flex gap-2 flex-wrap">
				<CommandButton title="Jump command" name="jump" />
				<CommandButton title="Color command" name="dudes color" />
				<CommandButton title="Grow command" name="dudes grow" />
				<CommandButton title="Sprite command" name="dudes sprite" />
				<CommandButton title="Leave command" name="dudes leave" />
			</div>
			<div class="mt-4">
				<div class="flex items-center justify-between mb-2">
					<span class="text-sm font-medium">{{ t('overlays.chat.presets') }}</span>
					<UiButton
						v-if="addable"
						size="sm"
						variant="outline"
						@click="handleAdd"
					>
						Add Preset
					</UiButton>
				</div>
				<UiTabs v-if="entities?.dudesGetAll?.length" v-model="openedTab">
					<UiTabsList class="w-full">
						<UiTabsTrigger
							v-for="(entity, entityIndex) in entities?.dudesGetAll"
							:key="entity.id"
							:value="entity.id!"
							class="relative group"
						>
							#{{ entityIndex + 1 }}
							<UiButton
								v-if="userCanEditOverlays"
								size="icon"
								variant="ghost"
								class="h-4 w-4 ml-1 opacity-0 group-hover:opacity-100"
								@click.stop="handleClose(entity.id!)"
							>
								<span class="text-xs">×</span>
							</UiButton>
						</UiTabsTrigger>
					</UiTabsList>
					<UiTabsContent
						v-for="entity in entities?.dudesGetAll"
						:key="entity.id"
						:value="entity.id!"
					>
						<DudesSettingsForm />
					</UiTabsContent>
				</UiTabs>
				<UiAlert v-else class="mt-2">
					<UiAlertDescription>
						Create new overlay for edit settings
					</UiAlertDescription>
				</UiAlert>
			</div>
		</div>
	</div>

	<UiAlertDialog v-model:open="deleteDialogOpen">
		<UiAlertDialogContent>
			<UiAlertDialogHeader>
				<UiAlertDialogTitle>Delete preset</UiAlertDialogTitle>
				<UiAlertDialogDescription>
					Are you sure you want to delete this preset?
				</UiAlertDialogDescription>
			</UiAlertDialogHeader>
			<UiAlertDialogFooter>
				<UiAlertDialogCancel>Cancel</UiAlertDialogCancel>
				<UiAlertDialogAction @click="confirmDelete">Delete</UiAlertDialogAction>
			</UiAlertDialogFooter>
		</UiAlertDialogContent>
	</UiAlertDialog>
</template>

<style scoped>
@import '../styles.css';

.iframe {
	height: 100%;
	width: 100%;
	aspect-ratio: 16 / 9;
	border: 0;
	margin-top: 8px;
	border: 1px solid hsl(var(--border));
	border-radius: 0.5rem;
}
</style>
