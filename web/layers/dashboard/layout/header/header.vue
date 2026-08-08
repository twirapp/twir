<script setup lang="ts">
import { useMediaQuery } from '@vueuse/core'
import { computed, ref } from 'vue'
import { useRealtimeDashboardStats } from '~~/layers/dashboard/api/dashboard'
import CommandMenu from '~~/layers/dashboard/components/command-menu/CommandMenu.vue'

import StreamInfoEditor from '../stream-info-editor.vue'
import { useHeaderLayout } from './composables/use-header-layout.js'
import HeaderBotStatus from './header-bot-status.vue'
import HeaderProfile from './header-profile.vue'
import StatsLayoutAggregate from './layouts/stats-layout-aggregate.vue'
import StatsLayoutGrid from './layouts/stats-layout-grid.vue'
import StatsLayoutRows from './layouts/stats-layout-rows.vue'
import LayoutSwitcher from './ui/layout-switcher.vue'

import { Platform } from '~/gql/graphql.js'

const { stats } = useRealtimeDashboardStats()

// Mobile detection
const isDesktop = useMediaQuery('(min-width: 768px)')

const { layout } = useHeaderLayout()

const activeLayoutComponent = computed(() => {
	switch (layout.value) {
		case 'rows':
			return StatsLayoutRows
		case 'grid':
			return StatsLayoutGrid
		default:
			return StatsLayoutAggregate
	}
})

const streamInfoEditorOpen = ref(false)
const editingPlatform = ref<Platform>(Platform.Twitch)

const editingPlatformStats = computed(() =>
	stats.value?.platforms.find((platform) => platform.platform === editingPlatform.value)
)

function openInfoEditor(platform: Platform) {
	editingPlatform.value = platform
	streamInfoEditorOpen.value = true
}
</script>

<template>
	<div
		class="bg-card border-b-border flex w-full flex-wrap justify-between gap-2 border-b px-2 py-1"
	>
		<div class="flex flex-col flex-wrap gap-2 py-1 md:flex-row">
			<!-- Mobile search icon -->
			<CommandMenu
				v-if="!isDesktop"
				:icon-only="true"
			/>

			<!-- Active stats layout -->
			<component
				:is="activeLayoutComponent"
				v-if="isDesktop"
				:stats="stats"
				@edit-stream-info="openInfoEditor"
			/>

			<!-- Layout switcher -->
			<LayoutSwitcher v-if="isDesktop" />
		</div>

		<div class="flex-end ml-auto flex flex-wrap items-center justify-end gap-2">
			<CommandMenu v-if="isDesktop" />
			<HeaderBotStatus />
			<HeaderProfile />
		</div>
	</div>

		<StreamInfoEditor
			v-model:open="streamInfoEditorOpen"
			:platform="editingPlatform"
			:title="editingPlatformStats?.title ?? undefined"
			:category-id="editingPlatformStats?.categoryId ?? undefined"
			:category-name="editingPlatformStats?.categoryName ?? undefined"
		/>
</template>
