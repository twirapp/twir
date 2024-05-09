<script setup lang="ts">
import {
	NAlert,
	NButton,
	NSpace,
	NTabPane,
	NTabs,
	useThemeVars,
} from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

import DudesSettingsForm from './dudes-settings-form.vue'
import { useDudesForm } from './use-dudes-form.js'
import { useDudesIframe } from './use-dudes-frame.js'

import type { DudesSettingsWithOptionalId } from './dudes-settings.js'

import {
	useDudesOverlayManager,
	useProfile,
	useUserAccessFlagChecker,
} from '@/api/index.js'
import { useNaiveDiscrete } from '@/composables/use-naive-discrete.js'
import CommandButton from '@/features/commands/components/command-button.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const themeVars = useThemeVars()
const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)
const dudesOverlayManager = useDudesOverlayManager()
const creator = dudesOverlayManager.useCreate()
const deleter = dudesOverlayManager.useDelete()
const { data: profile } = useProfile()

const { t } = useI18n()
const { dialog } = useNaiveDiscrete()

const {
	data: entities,
} = dudesOverlayManager.useGetAll()

const openedTab = ref<string>()

const { dudesIframe, sendIframeMessage } = useDudesIframe()
const dudesIframeUrl = computed(() => {
	if (!profile.value || !openedTab.value) return null
	return `${window.location.origin}/overlays/${profile.value.apiKey}/dudes?id=${openedTab.value}`
})

const { setData, getDefaultSettings } = useDudesForm()

function resetTab() {
	if (!entities.value?.settings.at(0)) {
		openedTab.value = undefined
		return
	}

	openedTab.value = entities.value.settings.at(0)?.id
}

watch(entities, () => {
	resetTab()
}, { immediate: true })

watch(openedTab, (v) => {
	const entity = entities.value?.settings.find(s => s.id === v) as DudesSettingsWithOptionalId
	if (!entity) return
	setData(entity)
})

async function handleClose(id: string) {
	dialog.create({
		title: 'Delete preset',
		content: 'Are you sure you want to delete this preset?',
		positiveText: 'Delete',
		negativeText: 'Cancel',
		showIcon: false,
		onPositiveClick: async () => {
			const entity = entities.value?.settings.find(s => s.id === id)
			if (!entity?.id) return

			await deleter.mutateAsync(entity.id)
			resetTab()
		},
	})
}

async function handleAdd() {
	await creator.mutateAsync(getDefaultSettings())
}

const addable = computed(() => {
	return userCanEditOverlays.value && (entities.value?.settings.length ?? 0) < 5
})
</script>

<template>
	<div class="flex gap-10" style="height: calc(100% - var(--layout-header-height));">
		<div class="relative w-[70%]">
			<div v-if="dudesIframeUrl">
				<iframe
					ref="dudesIframe"
					style="position: absolute;"
					:src="dudesIframeUrl"
					class="iframe"
				/>
				<NSpace :size="6" class="absolute top-[18px] left-2">
					<NButton @click="sendIframeMessage('spawn-emote')">
						Emote
					</NButton>
					<NButton @click="sendIframeMessage('show-message')">
						Message
					</NButton>
					<NButton @click="sendIframeMessage('jump')">
						Jump
					</NButton>
					<NButton @click="sendIframeMessage('grow')">
						Grow
					</NButton>
					<NButton @click="sendIframeMessage('leave')">
						Leave
					</NButton>
					<NButton @click="sendIframeMessage('reset')">
						Reset
					</NButton>
				</NSpace>
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
			<NTabs
				v-model:value="openedTab"
				type="card"
				:closable="userCanEditOverlays"
				:addable="addable"
				style="margin-top: 1rem;"
				tab-style="min-width: 80px;"
				@close="handleClose"
				@add="handleAdd"
			>
				<template #prefix>
					{{ t('overlays.chat.presets') }}
				</template>
				<template v-if="entities?.settings.length">
					<NTabPane
						v-for="(entity, entityIndex) in entities?.settings"
						:key="entity.id"
						:tab="`#${entityIndex + 1}`"
						:name="entity.id!"
					>
						<DudesSettingsForm />
					</NTabPane>
				</template>
			</NTabs>
			<NAlert v-if="!entities?.settings.length" type="info" class="mt-2">
				Create new overlay for edit settings
			</NAlert>
		</div>
	</div>
</template>

<style scope>
@import '../styles.css';

.iframe {
	height: 100%;
	width: 100%;
	aspect-ratio: 16 / 9;
	border: 0;
	margin-top: 8px;
	border: 1px solid v-bind('themeVars.borderColor');
	border-radius: 8px;
}
</style>
