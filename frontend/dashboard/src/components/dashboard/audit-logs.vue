<script setup lang="ts">
import { IconExternalLink } from '@tabler/icons-vue'
import { UseTimeAgo } from '@vueuse/components'
import { NButton, NScrollbar, NTooltip } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import type { BadgeVariants } from '@/components/ui/badge'

import { useProfile } from '@/api'
import { mapOperationTypeToTranslate, mapSystemToTranslate, useAuditLogs } from '@/api/audit-logs'
import Card from '@/components/dashboard/card.vue'
import { Badge } from '@/components/ui/badge'
import { AuditOperationType } from '@/gql/graphql'

const props = defineProps<{
	popup?: boolean
}>()

const { logs } = useAuditLogs()

const { t } = useI18n()

function computeValueHint(oldValue: string | null | undefined, newValue: string | null | undefined) {
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
		`height=${height},width=${width},top=${top},left=${left},status=0,location=0,menubar=0,toolbar=0`,
	)
}
</script>

<template>
	<Card :content-style="{ padding: '0px', height: '80%' }" :popup="props.popup">
		<template #header-extra>
			<NTooltip trigger="hover" placement="bottom">
				<template #trigger>
					<NButton size="small" text @click="openPopup">
						<IconExternalLink />
					</NButton>
				</template>

				{{ t('sharedButtons.popout') }}
			</NTooltip>
		</template>

		<NScrollbar trigger="none">
			<TransitionGroup name="list">
				<div v-for="log of logs" :key="log.id" class="flex flex-col md:flex-row md:justify-between p-2 pr-4 border-b-[color:var(--n-border-color)] border-b border-solid">
					<div class="flex h-full flex-col sm:flex-col md:flex-row min-h-[40px] gap-2.5 px-2.5 md:items-center">
						<div v-if="log.user" class="flex gap-2 items-center">
							<img class="size-4 rounded-full" :src="log.user.profileImageUrl" />
							<span>{{ log.user.displayName }}</span>
						</div>
						<Badge :variant="computeOperationBadgeVariant(log.operationType)">
							{{ t(mapOperationTypeToTranslate(log.operationType)).toLocaleLowerCase() }}
							<template v-if="log.system">
								{{ t(mapSystemToTranslate(log.system), {}, { default: log.system }).toLocaleLowerCase() }}
							</template>
						</Badge>
						<Badge v-if="computeValueHint(log.oldValue, log.newValue)" variant="secondary">
							{{ computeValueHint(log.oldValue, log.newValue) }}
						</Badge>
					</div>

					<Badge variant="outline" class="mt-2 md:mt-0 flex w-fit">
						<UseTimeAgo v-slot="{ timeAgo }" :time="new Date(log.createdAt)" show-second>
							{{ timeAgo }}
						</UseTimeAgo>
					</Badge>
				</div>
			</TransitionGroup>
		</NScrollbar>
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
