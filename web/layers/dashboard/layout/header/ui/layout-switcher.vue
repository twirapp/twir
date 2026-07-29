<script setup lang="ts">
import Popover from '@/components/ui/popover/Popover.vue'
import PopoverContent from '@/components/ui/popover/PopoverContent.vue'
import PopoverTrigger from '@/components/ui/popover/PopoverTrigger.vue'

import { headerStatsLayouts, useHeaderLayout } from '../composables/use-header-layout.js'

import type { HeaderStatsLayout } from '../composables/use-header-layout.js'

const { layout } = useHeaderLayout()

const { t } = useI18n()

const layoutIcons: Record<HeaderStatsLayout, string> = {
	rows: 'lucide:list',
	aggregate: 'lucide:layers',
	grid: 'lucide:layout-grid',
}
</script>

<template>
	<Popover>
		<PopoverTrigger as-child>
			<button
				class="text-muted-foreground hover:text-foreground flex items-center gap-1.5 rounded-lg border border-transparent px-3 py-2 text-xs transition-all hover:border-white/10 hover:bg-white/5"
				:title="t('dashboard.statsLayouts.label')"
			>
				<Icon
					name="lucide:layout-grid"
					:size="14"
				/>
			</button>
		</PopoverTrigger>
		<PopoverContent
			class="w-64 p-2"
			align="start"
		>
			<p class="text-muted-foreground px-2 py-1 text-xs font-semibold">
				{{ t('dashboard.statsLayouts.label') }}
			</p>
			<button
				v-for="option in headerStatsLayouts"
				:key="option"
				class="text-foreground flex w-full items-center gap-2 rounded px-2 py-1.5 text-left transition-colors hover:bg-white/10"
				@click="layout = option"
			>
				<Icon
					:name="layoutIcons[option]"
					:size="14"
					class="text-muted-foreground flex-none"
				/>
				<span class="flex min-w-0 flex-col">
					<span class="text-xs font-medium">{{ t(`dashboard.statsLayouts.${option}`) }}</span>
					<span class="text-muted-foreground text-[11px] leading-tight">
						{{ t(`dashboard.statsLayouts.${option}Description`) }}
					</span>
				</span>
				<Icon
					v-if="layout === option"
					name="lucide:check"
					:size="14"
					class="text-primary ml-auto flex-none"
				/>
			</button>
		</PopoverContent>
	</Popover>
</template>
