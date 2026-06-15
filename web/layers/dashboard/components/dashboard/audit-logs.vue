<script setup lang="ts">
import { UseTimeAgo } from '@vueuse/components'
import {
	mapOperationTypeToTranslate,
	mapSystemToTranslate,
	useAuditLogs,
} from '~~/layers/dashboard/api/audit-logs'
import { useProfile } from '~~/layers/dashboard/api/auth'
import AuditLogUser from '~~/layers/dashboard/components/dashboard/audit-log-user.vue'
import Card from '~~/layers/dashboard/components/dashboard/card.vue'

import { Badge, type BadgeVariants } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { CardContent } from '@/components/ui/card'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from '@/components/ui/tooltip'
import { AuditOperationType } from '~/gql/graphql.js'

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
	<Card
		:popup="props.popup"
		class="@container flex min-h-0 flex-col"
	>
		<template #header-extra>
			<TooltipProvider>
				<Tooltip>
					<TooltipTrigger as-child>
						<Button
							size="sm"
							variant="ghost"
							@click="openPopup"
						>
							<Icon
								name="lucide:external-link"
								class="h-4 w-4"
							/>
						</Button>
					</TooltipTrigger>
					<TooltipContent>
						{{ t('sharedButtons.popout') }}
					</TooltipContent>
				</Tooltip>
			</TooltipProvider>
		</template>

		<CardContent class="min-h-0 flex-1 overflow-hidden p-0">
			<ScrollArea class="h-full w-full">
				<TransitionGroup name="list">
					<div
						v-for="log of logs"
						:key="`${log.system}-${log.objectId}-${log.objectId}-${log.createdAt}`"
						class="border-b-border flex flex-col border-b border-solid p-1 pr-4 @lg:flex-row @lg:items-center @lg:justify-between"
					>
						<div class="flex h-full min-h-10 flex-row flex-wrap items-center gap-2.5 px-2.5 py-1">
							<Badge
								variant="outline"
								class="mt-2 flex h-fit w-fit shrink-0 bg-zinc-800 p-1 px-2 @lg:mt-0"
							>
								<UseTimeAgo
									v-slot="{ timeAgo }"
									:time="new Date(log.createdAt)"
									show-second
								>
									{{ timeAgo }}
								</UseTimeAgo>
							</Badge>
							<AuditLogUser
								v-if="log.user"
								:user="log.user"
								:platform="log.platform"
							/>
							<Badge
								:variant="computeOperationBadgeVariant(log.operationType)"
								class="shrink-0"
							>
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
							<Badge
								v-if="computeValueHint(log.oldValue, log.newValue)"
								variant="secondary"
								class="max-w-[200px] shrink-0 truncate"
							>
								{{ computeValueHint(log.oldValue, log.newValue) }}
							</Badge>
						</div>
					</div>
				</TransitionGroup>
			</ScrollArea>
		</CardContent>
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
