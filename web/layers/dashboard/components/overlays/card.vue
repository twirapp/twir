<script setup lang="ts">

import { useCopyOverlayLink } from './copyOverlayLink.js'

import { useUserAccessFlagChecker } from '~~/layers/dashboard/api/auth'
import Card from '~~/layers/dashboard/components/card/card.vue'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { ChannelRolePermissionEnum } from '~/gql/graphql.js'

const props = withDefaults(
	defineProps<{
		description: string
		title: string
		overlayPath?: string
		icon: string
		iconStroke?: number
		showSettings?: boolean
		copyDisabled?: boolean
		showCopy?: boolean
	}>(),
	{
		showSettings: true,
		copyDisabled: false,
		showCopy: true,
		iconStroke: 1,
		overlayPath: '',
	}
)

defineEmits<{
	openSettings: []
}>()

const { t } = useI18n()
const { copyOverlayLink } = useCopyOverlayLink(props.overlayPath)

const userCanEditOverlays = useUserAccessFlagChecker(ChannelRolePermissionEnum.ManageOverlays)
</script>

<template>
	<Card :title="title" :icon="icon" :icon-stroke="iconStroke" class="h-full">
		<template #content>
			{{ description }}
		</template>

		<template #footer>
			<Button
				v-if="showSettings"
				:disabled="!userCanEditOverlays"
				variant="outline"
				class="flex items-center gap-2"
				@click="$emit('openSettings')"
			>
				<Icon name="lucide:settings" class="size-4" />
				<span>{{ t('sharedButtons.settings') }}</span>
			</Button>
			<TooltipProvider v-if="showCopy">
				<Tooltip>
					<TooltipTrigger as-child>
						<Button
							:disabled="copyDisabled"
							class="flex items-center gap-2"
							@click="copyOverlayLink()"
						>
							<Icon name="lucide:copy" class="size-4" />
							<span>{{ t('overlays.copyOverlayLink') }}</span>
						</Button>
					</TooltipTrigger>
					<TooltipContent>
						{{ t('overlays.uncongirured') }}
					</TooltipContent>
				</Tooltip>
			</TooltipProvider>
		</template>
	</Card>
</template>
