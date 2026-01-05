<script setup lang="ts">
import { ExternalLinkIcon } from 'lucide-vue-next'
import { UseTimeAgo } from '@vueuse/components'
import { useI18n } from 'vue-i18n'

import { useProfile } from '@/api'
import { mapOperationTypeToTranslate, mapSystemToTranslate, useAuditLogs } from '@/api/audit-logs'
import Card from '@/components/dashboard/card.vue'
import { Badge, type BadgeVariants } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { AuditOperationType } from '@/gql/graphql'

const props = defineProps<{
	popup?: boolean
}>()

const { logs } = useAuditLogs()

const { t } = useI18n()

function computeValueHint(
	oldValue: string | null | undefined,
	newValue: string | null | undefined
) {
	if (!oldValue && !newValue) return null
	const oldJson = oldValue ? JSON.parse(oldValue) : null
	const newJson = newValue ? JSON.parse(newValue) : null

	const value = newJson || oldJson
	if (!value) return null

	return value.name ?? value.text ?? value.Name ?? null
}

function computeOperationBadgeVariant(operation: AuditOperationType): BadgeVariants['variant'] {
	switch (operation) {
		case AuditOperationType.Create:
			return 'success'
		case AuditOperationType.Delete:
			return 'destructive'
		case AuditOperationType.Update:
			return 'default'
	}
}

const { data: profile } = useProfile()

function openPopup() {
	if (!profile.value) return

	const height = 800
	const width = 500
	const top = Math.max(0, (screen.height - height) / 2)
	const left = Math.max(0, (screen.width - width) / 2)

	window.open(
		`${window.location.origin}/dashboard/popup/widgets/audit-log?apiKey=${profile.value.apiKey}`,
		'_blank',
		`height=${height},width=${width},top=${top},left=${left},status=0,location=0,menubar=0,toolbar=0`
	)
}
</script>

<template>
	<Card :content-style="{ padding: '0px', height: '80%' }" :popup="props.popup" class="@container">
		<template #header-extra>
			<TooltipProvider>
				<Tooltip>
					<TooltipTrigger as-child>
						<Button size="sm" variant="ghost" @click="openPopup">
							<ExternalLinkIcon class="h-4 w-4" />
						</Button>
					</TooltipTrigger>
					<TooltipContent>
						{{ t('sharedButtons.popout') }}
					</TooltipContent>
				</Tooltip>
			</TooltipProvider>
		</template>

		<ScrollArea class="h-full">
			<TransitionGroup name="list">
				<div
					v-for="log of logs"
					:key="`${log.system}-${log.objectId}-${log.objectId}-${log.createdAt}`"
					class="flex flex-col @lg:flex-row @lg:justify-between p-1 pr-4 border-b-border border-b border-solid @lg:items-center"
				>
					<div
						class="flex h-full flex-row items-center min-h-10 gap-2.5 px-2.5 py-1"
					>
						<Badge variant="outline" class="mt-2 @lg:mt-0 flex w-fit p-1 h-fit px-2 bg-zinc-800">
							<UseTimeAgo v-slot="{ timeAgo }" :time="new Date(log.createdAt)" show-second>
								{{ timeAgo }}
							</UseTimeAgo>
						</Badge>
						<div v-if="log.user" class="flex gap-2 items-center">
							<img class="size-4 rounded-full" :src="log.user.profileImageUrl" />
							<span>{{ log.user.displayName }}</span>
						</div>
						<Badge :variant="computeOperationBadgeVariant(log.operationType)">
							{{ t(mapOperationTypeToTranslate(log.operationType)).toLocaleLowerCase() }}
							<template v-if="log.system">
								{{
									t(
										mapSystemToTranslate(log.system),
										{},
										{ default: log.system }
									).toLocaleLowerCase()
								}}
							</template>
						</Badge>
						<Badge v-if="computeValueHint(log.oldValue, log.newValue)" variant="secondary">
							{{ computeValueHint(log.oldValue, log.newValue) }}
						</Badge>
					</div>
				</div>
			</TransitionGroup>
		</ScrollArea>
	</Card>
</template>

<style scoped>
.list-enter-active,
.list-leave-active {
	transition: all 0.5s ease;
}

.list-enter-from,
.list-leave-to {
	opacity: 0;
	transform: translateY(30px);
}
</style>
