<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import DudesSettingsForm from './dudes-settings-form.vue'
import { useDudesForm } from './use-dudes-form.js'
import { useDudesIframe } from './use-dudes-frame.js'

import type { DudesSettingsWithOptionalId } from './dudes-settings.js'

import { useProfile, useUserAccessFlagChecker } from '@/api/auth'
import { useDudesOverlayManager } from '@/api/overlays/dudes'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
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
import CommandButton from '@/features/commands/ui/command-button.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

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
					<Button size="sm" @click="sendIframeMessage('spawn-emote')"> Emote </Button>
					<Button size="sm" @click="sendIframeMessage('show-message')"> Message </Button>
					<Button size="sm" @click="sendIframeMessage('jump')"> Jump </Button>
					<Button size="sm" @click="sendIframeMessage('grow')"> Grow </Button>
					<Button size="sm" @click="sendIframeMessage('leave')"> Leave </Button>
					<Button size="sm" @click="sendIframeMessage('reset')"> Reset </Button>
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
					<Button
						v-if="addable"
						size="sm"
						variant="outline"
						@click="handleAdd"
					>
						Add Preset
					</Button>
				</div>
				<Tabs v-if="entities?.dudesGetAll?.length" v-model="openedTab">
					<TabsList class="w-full">
						<TabsTrigger
							v-for="(entity, entityIndex) in entities?.dudesGetAll"
							:key="entity.id"
							:value="entity.id!"
							class="relative group"
						>
							#{{ entityIndex + 1 }}
							<Button
								v-if="userCanEditOverlays"
								size="icon"
								variant="ghost"
								class="h-4 w-4 ml-1 opacity-0 group-hover:opacity-100"
								@click.stop="handleClose(entity.id!)"
							>
								<span class="text-xs">×</span>
							</Button>
						</TabsTrigger>
					</TabsList>
					<TabsContent
						v-for="entity in entities?.dudesGetAll"
						:key="entity.id"
						:value="entity.id!"
					>
						<DudesSettingsForm />
					</TabsContent>
				</Tabs>
				<Alert v-else class="mt-2">
					<AlertDescription>
						Create new overlay for edit settings
					</AlertDescription>
				</Alert>
			</div>
		</div>
	</div>

	<AlertDialog v-model:open="deleteDialogOpen">
		<AlertDialogContent>
			<AlertDialogHeader>
				<AlertDialogTitle>Delete preset</AlertDialogTitle>
				<AlertDialogDescription>
					Are you sure you want to delete this preset?
				</AlertDialogDescription>
			</AlertDialogHeader>
			<AlertDialogFooter>
				<AlertDialogCancel>Cancel</AlertDialogCancel>
				<AlertDialogAction @click="confirmDelete">Delete</AlertDialogAction>
			</AlertDialogFooter>
		</AlertDialogContent>
	</AlertDialog>
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
