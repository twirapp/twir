<script setup lang="ts">
import { usePlatformStats } from '../composables/use-platform-stats.js'
import PlatformIcon from '../ui/platform-icon.vue'

import { Platform } from '~/gql/graphql.js'

import type { DashboardPlatformStats, DashboardStats } from '../composables/use-platform-stats.js'

interface Props {
	stats: DashboardStats | undefined
}

const props = defineProps<Props>()

const emit = defineEmits<{
	editStreamInfo: [platform: Platform]
}>()

const { t } = useI18n()

const { sortedPlatforms, uptimes } = usePlatformStats(() => props.stats)

function canEdit(platform: DashboardPlatformStats): boolean {
	return (
		platform.canEditInfo &&
		(platform.platform === Platform.Twitch || platform.platform === Platform.Kick)
	)
}
</script>

<template>
	<!-- Platform strips only, all in one horizontal row; no global stats widgets in this layout -->
	<div class="flex items-center gap-2 overflow-x-auto py-0.5">
		<div
			v-for="platform in sortedPlatforms"
			:key="platform.platform"
			class="header-widget"
		>
			<div class="flex items-center gap-2 whitespace-nowrap">
				<PlatformIcon
					:platform="platform.platform"
					size="sm"
					:live="platform.isLive"
				/>

				<span
					v-if="!platform.isLive"
					class="bg-muted text-muted-foreground flex-none rounded px-1.5 py-0.5 text-[10px] font-bold"
				>
					{{ t('dashboard.platformStats.offline') }}
				</span>

				<span class="text-foreground min-w-0 max-w-40 flex-1 truncate text-[13px] leading-tight font-semibold">
					{{ platform.title || t('dashboard.platformStats.noTitle') }}
				</span>
				<button
					v-if="canEdit(platform)"
					class="text-muted-foreground hover:text-foreground flex flex-none items-center transition-colors"
					:title="t('dashboard.statsWidgets.streamInfo.editStreamInfo')"
					@click.stop="emit('editStreamInfo', platform.platform)"
				>
					<Icon
						name="lucide:pencil"
						:size="11"
					/>
				</button>

				<span
					class="text-muted-foreground max-w-28 flex-none truncate rounded-full border border-white/10 bg-white/5 px-2 py-0.5 text-[11px]"
				>
					{{ platform.categoryName || t('dashboard.platformStats.noCategory') }}
				</span>

				<span class="h-4 w-px flex-none bg-white/10" />

				<span class="flex flex-none items-center gap-1 text-[12px] font-semibold tabular-nums">
					<Icon
						name="lucide:clock"
						:size="12"
						class="text-muted-foreground"
					/>
					{{ platform.isLive ? (uptimes[platform.platform] ?? '—') : '—' }}
				</span>
				<span class="flex flex-none items-center gap-1 text-[12px] font-semibold tabular-nums">
					<Icon
						name="lucide:users"
						:size="12"
						class="text-muted-foreground"
					/>
					{{ platform.isLive ? (platform.viewers ?? 0) : '—' }}
				</span>
				<span class="flex flex-none items-center gap-1 text-[12px] font-semibold tabular-nums">
					<Icon
						name="lucide:user-plus"
						:size="12"
						class="text-muted-foreground"
					/>
					{{ platform.followers ?? '—' }}
				</span>
			</div>
		</div>
	</div>
</template>
