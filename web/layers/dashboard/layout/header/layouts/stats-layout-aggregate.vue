<script setup lang="ts">
import { ref } from 'vue'

import Popover from '@/components/ui/popover/Popover.vue'
import PopoverContent from '@/components/ui/popover/PopoverContent.vue'
import PopoverTrigger from '@/components/ui/popover/PopoverTrigger.vue'

import { usePlatformStats } from '../composables/use-platform-stats.js'
import GlobalStatsWidgets from '../global-stats-widgets.vue'
import PlatformIcon from '@/components/platform/platform-icon.vue'
import PlatformStatsRows from '../ui/platform-stats-rows.vue'
import PlatformTitleRow from '../ui/platform-title-row.vue'
import StatWidget from '../ui/stat-widget.vue'

import type { Platform } from '~/gql/graphql.js'

import type { DashboardStats } from '../composables/use-platform-stats.js'

interface Props {
	stats: DashboardStats | undefined
}

const props = defineProps<Props>()

const emit = defineEmits<{
	editStreamInfo: [platform: Platform]
}>()

const { t } = useI18n()

const { sortedPlatforms, totalViewers, totalFollowers, primaryPlatform, uptimes, globalUptime } =
	usePlatformStats(() => props.stats)

const breakdownOpen = ref(false)

function onTitleWidgetClick() {
	const primary = primaryPlatform.value
	if (primary?.canEditInfo) {
		emit('editStreamInfo', primary.platform)
		return
	}
	breakdownOpen.value = true
}
</script>

<template>
	<!-- Title widget: primary platform, click opens Twitch editor or the breakdown -->
	<StatWidget
		interactive
		class="header-widget-stream"
		@click="onTitleWidgetClick"
	>
		<div class="header-widget-content">
			<div class="flex items-center gap-2">
				<PlatformIcon
					v-if="primaryPlatform"
					:platform="primaryPlatform.platform"
					:live="primaryPlatform.isLive"
				/>
				<div class="flex min-w-0 flex-1 flex-col">
					<p class="header-widget-value truncate">
						{{ primaryPlatform?.title || t('dashboard.platformStats.noTitle') }}
					</p>
					<p class="header-widget-label truncate">
						{{ primaryPlatform?.categoryName || t('dashboard.platformStats.noCategory') }}
					</p>
				</div>
				<Icon
				v-if="primaryPlatform?.canEditInfo"
					name="tabler:edit"
					class="text-muted-foreground h-3.5 w-3.5 flex-shrink-0"
				/>
				<Icon
					v-else
					name="lucide:chevron-down"
					class="text-muted-foreground h-3.5 w-3.5 flex-shrink-0"
				/>
			</div>
		</div>
	</StatWidget>

	<!-- Global widgets; viewers/followers show totals + stacked icons + breakdown popover -->
	<GlobalStatsWidgets
		:stats="stats"
		:uptime="globalUptime"
	>
		<template #viewers>
			<Popover v-model:open="breakdownOpen">
				<PopoverTrigger as-child>
					<div class="header-widget-content cursor-pointer">
						<div class="flex items-center gap-1.5">
							<p class="header-widget-value">
								{{ totalViewers }}
							</p>
							<span class="flex flex-none items-center">
								<PlatformIcon
									v-for="(platform, index) in sortedPlatforms"
									:key="platform.platform"
									:platform="platform.platform"
									size="sm"
									:class="{ '-ml-1': index > 0 }"
								/>
							</span>
							<Icon
								name="lucide:chevron-down"
								class="text-muted-foreground h-3 w-3 flex-none"
							/>
						</div>
						<p class="header-widget-label">
							{{ t('dashboard.statsWidgets.viewers') }}
						</p>
					</div>
				</PopoverTrigger>
				<PopoverContent
					class="w-[26rem] p-2"
					align="start"
				>
					<p class="text-muted-foreground px-2 py-1 text-xs font-semibold">
						{{ t('dashboard.platformStats.breakdown') }}
					</p>
					<PlatformStatsRows
						mode="table"
						:platforms="sortedPlatforms"
						:uptimes="uptimes"
					/>
					<div class="my-1 h-px bg-white/10" />
					<p class="text-muted-foreground px-2 py-1 text-xs font-semibold">
						{{ t('dashboard.platformStats.streamTitles') }}
					</p>
					<PlatformTitleRow
						v-for="platform in sortedPlatforms"
						:key="platform.platform"
						:platform-stat="platform"
						@edit="emit('editStreamInfo', $event)"
					/>
				</PopoverContent>
			</Popover>
		</template>

		<template #followers>
			<Popover>
				<PopoverTrigger as-child>
					<div class="header-widget-content cursor-pointer">
						<div class="flex items-center gap-1.5">
							<p class="header-widget-value">
								{{ totalFollowers }}
							</p>
							<span class="flex flex-none items-center">
								<PlatformIcon
									v-for="(platform, index) in sortedPlatforms"
									:key="platform.platform"
									:platform="platform.platform"
									size="sm"
									:class="{ '-ml-1': index > 0 }"
								/>
							</span>
							<Icon
								name="lucide:chevron-down"
								class="text-muted-foreground h-3 w-3 flex-none"
							/>
						</div>
						<p class="header-widget-label">
							{{ t('dashboard.statsWidgets.followers') }}
						</p>
					</div>
				</PopoverTrigger>
				<PopoverContent
					class="w-[26rem] p-2"
					align="start"
				>
					<p class="text-muted-foreground px-2 py-1 text-xs font-semibold">
						{{ t('dashboard.platformStats.breakdown') }}
					</p>
					<PlatformStatsRows
						mode="table"
						:platforms="sortedPlatforms"
						:uptimes="uptimes"
					/>
					<div class="my-1 h-px bg-white/10" />
					<p class="text-muted-foreground px-2 py-1 text-xs font-semibold">
						{{ t('dashboard.platformStats.streamTitles') }}
					</p>
					<PlatformTitleRow
						v-for="platform in sortedPlatforms"
						:key="platform.platform"
						:platform-stat="platform"
						@edit="emit('editStreamInfo', $event)"
					/>
				</PopoverContent>
			</Popover>
		</template>
	</GlobalStatsWidgets>
</template>
