<script setup lang="ts">
import PlatformIcon from './platform-icon.vue'

import { Platform } from '~/gql/graphql.js'

import type { DashboardPlatformStats } from '../composables/use-platform-stats.js'

interface Props {
	platforms: DashboardPlatformStats[]
	mode?: 'compact' | 'table'
	field?: 'viewers' | 'followers'
	uptimes?: Partial<Record<Platform, string>>
}

const props = withDefaults(defineProps<Props>(), {
	mode: 'compact',
	field: 'viewers',
	uptimes: () => ({}),
})

const { t } = useI18n()

const platformNames: Record<Platform, string> = {
	[Platform.Twitch]: 'Twitch',
	[Platform.Kick]: 'Kick',
	[Platform.VkVideoLive]: 'VK',
	[Platform.Youtube]: 'YouTube',
}

function compactValue(platform: DashboardPlatformStats): string | number {
	if (props.field === 'viewers') {
		if (!platform.isLive) return '—'
		return platform.viewers ?? 0
	}
	return platform.followers ?? '—'
}

function tableViewers(platform: DashboardPlatformStats): string | number {
	if (!platform.isLive) return '—'
	return platform.viewers ?? 0
}
</script>

<template>
	<!-- Compact mode: per-platform icon + number in one horizontal row (rows-layout widgets) -->
	<div
		v-if="mode === 'compact'"
		class="flex items-center gap-2.5"
	>
		<span
			v-for="platform in platforms"
			:key="platform.platform"
			class="flex items-center gap-1 text-sm leading-tight font-semibold tabular-nums"
			:class="{ 'opacity-40': !platform.isLive || compactValue(platform) === '—' }"
		>
			<PlatformIcon
				:platform="platform.platform"
				size="sm"
			/>
			{{ compactValue(platform) }}
		</span>
	</div>

	<!-- Table mode: Platform / Viewers / Followers / Uptime breakdown (aggregate popover) -->
	<div
		v-else
		class="flex flex-col"
	>
		<div
			class="text-muted-foreground grid grid-cols-[1fr_64px_72px_72px] items-center gap-2 px-2 pt-1 pb-1.5 text-[10px] font-bold tracking-wider uppercase"
		>
			<span>{{ t('dashboard.platformStats.platform') }}</span>
			<span class="text-right">{{ t('dashboard.statsWidgets.viewers') }}</span>
			<span class="text-right">{{ t('dashboard.statsWidgets.followers') }}</span>
			<span class="text-right">{{ t('dashboard.statsWidgets.uptime') }}</span>
		</div>
		<div
			v-for="platform in platforms"
			:key="platform.platform"
			class="grid grid-cols-[1fr_64px_72px_72px] items-center gap-2 rounded-md px-2 py-1.5 text-[13px] hover:bg-white/5"
			:class="{ 'opacity-50': !platform.isLive }"
		>
			<span class="flex min-w-0 items-center gap-2 font-semibold">
				<PlatformIcon
					:platform="platform.platform"
					:live="platform.isLive"
				/>
				<span class="truncate">{{ platformNames[platform.platform] }}</span>
				<span
					v-if="!platform.isLive"
					class="bg-muted text-muted-foreground flex-none rounded px-1.5 py-0.5 text-[10px] font-bold"
				>
					{{ t('dashboard.platformStats.offline') }}
				</span>
			</span>
			<span class="text-right font-semibold tabular-nums">{{ tableViewers(platform) }}</span>
			<span class="text-right font-semibold tabular-nums">{{ platform.followers ?? '—' }}</span>
			<span class="text-right font-semibold tabular-nums">{{ uptimes[platform.platform] ?? '—' }}</span>
		</div>
	</div>
</template>
