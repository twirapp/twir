<script setup lang="ts">
import Popover from '@/components/ui/popover/Popover.vue'
import PopoverContent from '@/components/ui/popover/PopoverContent.vue'
import PopoverTrigger from '@/components/ui/popover/PopoverTrigger.vue'

import { usePlatformStats } from '../composables/use-platform-stats.js'
import GlobalStatsWidgets from '../global-stats-widgets.vue'
import PlatformIcon from '../ui/platform-icon.vue'
import PlatformStatsRows from '../ui/platform-stats-rows.vue'
import PlatformTitleRow from '../ui/platform-title-row.vue'
import StatWidget from '../ui/stat-widget.vue'

import type { DashboardStats } from '../composables/use-platform-stats.js'
import type { Platform } from '~/gql/graphql.js'

interface Props {
	stats: DashboardStats | undefined
}

const props = defineProps<Props>()

const emit = defineEmits<{
	editStreamInfo: [platform: Platform]
}>()

const { t } = useI18n()

const { platforms, sortedPlatforms, livePlatformsCount, uptimes, globalUptime } =
	usePlatformStats(() => props.stats)
</script>

<template>
	<!-- Stream titles widget: cluster of platform icons + per-platform rows in a popover -->
	<Popover>
		<PopoverTrigger as-child>
			<StatWidget
				interactive
				class="header-widget-stream"
			>
				<div class="header-widget-content">
					<div class="flex items-center gap-2">
						<span class="flex flex-none items-center">
							<PlatformIcon
								v-for="(platform, index) in sortedPlatforms"
								:key="platform.platform"
								:platform="platform.platform"
								:class="{ '-ml-1': index > 0 }"
							/>
						</span>
						<div class="flex min-w-0 flex-1 flex-col">
							<p class="header-widget-value truncate">
								{{ t('dashboard.platformStats.streamTitles') }}
							</p>
							<p class="header-widget-label truncate">
								{{
									t('dashboard.platformStats.livePlatforms', {
										live: livePlatformsCount,
										total: platforms.length,
									})
								}}
							</p>
						</div>
						<Icon
							name="lucide:chevron-down"
							class="text-muted-foreground h-3.5 w-3.5 flex-shrink-0"
						/>
					</div>
				</div>
			</StatWidget>
		</PopoverTrigger>
		<PopoverContent
			class="w-96 p-2"
			align="start"
		>
			<p class="text-muted-foreground px-2 py-1 text-xs font-semibold">
				{{ t('dashboard.platformStats.streamTitles') }}
			</p>
			<PlatformTitleRow
				v-for="platform in sortedPlatforms"
				:key="platform.platform"
				:platform-stat="platform"
				:uptime="platform.isLive ? uptimes[platform.platform] : undefined"
				@edit="emit('editStreamInfo', $event)"
			/>
		</PopoverContent>
	</Popover>

	<!-- Global widgets; viewers/followers render per-platform rows -->
	<GlobalStatsWidgets
		:stats="stats"
		:uptime="globalUptime"
	>
		<template #viewers>
			<div class="header-widget-content">
				<PlatformStatsRows
					:platforms="sortedPlatforms"
					field="viewers"
				/>
				<p class="header-widget-label">
					{{ t('dashboard.statsWidgets.viewers') }}
				</p>
			</div>
		</template>
		<template #followers>
			<div class="header-widget-content">
				<PlatformStatsRows
					:platforms="sortedPlatforms"
					field="followers"
				/>
				<p class="header-widget-label">
					{{ t('dashboard.statsWidgets.followers') }}
				</p>
			</div>
		</template>
	</GlobalStatsWidgets>
</template>
