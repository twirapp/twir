<script setup lang="ts">
import { IconCopy, IconSettings } from '@tabler/icons-vue'
import { NButton, NTooltip } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import { useCopyOverlayLink } from './copyOverlayLink.js'

import type { FunctionalComponent } from 'vue'

import { useProfile, useUserAccessFlagChecker } from '@/api/index.js'
import Card from '@/components/card/card.vue'
import { ChannelRolePermissionEnum } from '@/gql/graphql'

const props = withDefaults(defineProps<{
	description: string
	title: string
	overlayPath?: string
	icon: FunctionalComponent
	iconStroke?: number
	showSettings?: boolean
	copyDisabled?: boolean
	showCopy?: boolean
}>(), {
	showSettings: true,
	copyDisabled: false,
	showCopy: true,
	iconStroke: 1,
	overlayPath: '',
})

defineEmits<{
	openSettings: []
}>()

const { t } = useI18n()
const { data: profile } = useProfile()

const { copyOverlayLink } = useCopyOverlayLink(props.overlayPath)

const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)
</script>

<template>
	<Card :title="title" :icon="icon" :icon-stroke="iconStroke" class="h-full">
		<template #content>
			{{ description }}
		</template>

		<template #footer>
			<NButton
				v-if="showSettings" :disabled="!userCanEditOverlays" secondary size="large"
				@click="$emit('openSettings')"
			>
				<div class="flex gap-1">
					<span>{{ t('sharedButtons.settings') }}</span>
					<IconSettings />
				</div>
			</NButton>
			<NTooltip v-if="showCopy" :disabled="profile?.id !== profile?.selectedDashboardId">
				<template #trigger>
					<NButton
						size="large"
						:disabled="copyDisabled || profile?.id !== profile?.selectedDashboardId"
						@click="copyOverlayLink()"
					>
						<div class="flex gap-1">
							<span>{{ t('overlays.copyOverlayLink') }}</span>
							<IconCopy />
						</div>
					</NButton>
				</template>
				<span v-if="profile?.id !== profile?.selectedDashboardId">{{ t('overlays.noAccess') }}</span>
				<span v-else>{{ t('overlays.uncongirured') }}</span>
			</NTooltip>
		</template>
	</Card>
</template>
