<script setup lang="ts">
import { computed } from 'vue'

import PlatformIcon from '@/components/platform/platform-icon.vue'

import type { Platform } from '~/gql/graphql.js'

import type { DashboardPlatformStats } from '../composables/use-platform-stats.js'

interface Props {
	platformStat: DashboardPlatformStats
	uptime?: string
	showIcon?: boolean
}

const props = withDefaults(defineProps<Props>(), {
	uptime: undefined,
	showIcon: true,
})

const emit = defineEmits<{
	edit: [platform: Platform]
}>()

const { t } = useI18n()

const editable = computed(() => props.platformStat.canEditInfo)
</script>

<template>
	<div
		class="flex items-center gap-2 rounded-md px-2 py-1.5"
		:class="{
			'opacity-50': !platformStat.isLive,
			'cursor-pointer transition-colors hover:bg-white/5': editable,
		}"
		@click="editable ? emit('edit', platformStat.platform) : undefined"
	>
		<PlatformIcon
			v-if="showIcon"
			:platform="platformStat.platform"
			:live="platformStat.isLive"
		/>
		<div class="flex min-w-0 flex-1 flex-col">
			<p class="text-foreground truncate text-[13px] leading-tight font-semibold">
				{{ platformStat.title || t('dashboard.platformStats.noTitle') }}
			</p>
			<p class="text-muted-foreground truncate text-[11px] leading-tight">
				{{ platformStat.categoryName || t('dashboard.platformStats.noCategory') }}
				<template v-if="uptime"> · {{ uptime }}</template>
			</p>
		</div>
		<button
			:disabled="!editable"
			class="text-muted-foreground hover:text-foreground hover:bg-muted flex h-6 w-6 flex-none items-center justify-center rounded-md transition-colors"
			:title="t(editable ? 'dashboard.statsWidgets.streamInfo.editStreamInfo' : 'dashboard.platformStats.editUnavailable')"
			@click.stop="editable ? emit('edit', platformStat.platform) : undefined"
		>
			<Icon
				name="lucide:pencil"
				:size="12"
			/>
		</button>
	</div>
</template>
