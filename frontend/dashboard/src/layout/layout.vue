<script setup lang="ts">
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'
import {
	NConfigProvider,
	NDialogProvider,
	NLayout,
	NLayoutContent,
	NLayoutHeader,
	NLayoutSider,
	NMessageProvider,
	NNotificationProvider,
	NScrollbar,
	NSpin,
	darkTheme,
	lightTheme,
} from 'naive-ui'
import { computed, ref, watch } from 'vue'
import { RouterView, useRouter } from 'vue-router'

import { Toaster as Sonner } from '@/components/ui/sonner'
import { Toaster } from '@/components/ui/toast'
import { useLayout } from '@/composables/use-layout'
import { useTheme } from '@/composables/use-theme.js'
import Header from '@/layout/header.vue'
import Sidebar from '@/layout/sidebar.vue'
import { useSidebarCollapseStore } from '@/layout/use-sidebar-collapse'

const { theme } = useTheme()
const themeStyles = computed(() => theme.value === 'dark' ? darkTheme : lightTheme)

const { layoutRef } = useLayout()
const isRouterReady = ref(false)
const router = useRouter()

router.isReady().finally(() => isRouterReady.value = true)

const breakPoints = useBreakpoints(breakpointsTailwind)
// If we are on a smaller than or equal to lg, we want the sidebar to collapse.
const smallerOrEqualLg = breakPoints.smallerOrEqual('lg')
// If we are on a smaller than or equal to md, we want the sidebar to hide and show hamburger menu with drawer.
const smallerOrEqualMd = breakPoints.smallerOrEqual('md')

const { isCollapsed, set: setIsCollapsed, toggle: toggleSidebar } = useSidebarCollapseStore()

watch(smallerOrEqualLg, (v) => {
	setIsCollapsed(v)
})
</script>

<template>
	<NConfigProvider
		:theme="themeStyles"
		class="h-full"
		:breakpoints="{ 'xs': 0, 's': 640, 'm': 1024, 'l': 1280, 'xl': 1536, 'xxl': 1920, '2xl': 2560 }"
	>
		<NNotificationProvider :max="5">
			<NMessageProvider :duration="2500" :closable="true">
				<NDialogProvider>
					<NLayout class="h-full">
						<NLayoutHeader bordered class="w-full" style="height: var(--layout-header-height)">
							<Header />
						</NLayoutHeader>

						<NLayout has-sider style="height: calc(100vh - var(--layout-header-height))">
							<NLayoutSider
								v-if="!smallerOrEqualMd"
								bordered
								collapse-mode="width"
								:collapsed-width="64"
								:width="240"
								show-trigger="arrow-circle"
								:native-scrollbar="false"
								:collapsed="isCollapsed"
								:show-collapsed-content="false"
								@update-collapsed="toggleSidebar"
							>
								<Sidebar />
							</NLayoutSider>
							<NLayoutContent ref="layoutRef">
								<div v-if="!isRouterReady" class="app-loader">
									<NSpin size="large" />
								</div>
								<RouterView v-else v-slot="{ Component, route }">
									<NScrollbar trigger="none">
										<transition :name="route.meta.transition as string || 'router'" mode="out-in">
											<div
												:key="route.path"
												:style="{
													padding: route.meta?.noPadding ? undefined : '24px',
													height: route.meta?.fullScreen ? 'calc(100% - var(--layout-header-height))' : 'auto',
												}"
											>
												<component :is="Component" />
											</div>
										</transition>
									</NScrollbar>
								</RouterView>

								<Toaster />
								<Sonner />
							</NLayoutContent>
						</NLayout>
					</NLayout>
				</NDialogProvider>
			</NMessageProvider>
		</NNotificationProvider>
	</NConfigProvider>
</template>

<style>
.router-enter-active,
.router-leave-active {
	transition: all 0.2s cubic-bezier(0, 0, 0.2, 1);
}

.router-enter-from,
.router-leave-to {
	opacity: 0;
	transform: scale(0.98);
}

.app-loader {
	display: flex;
	justify-content: center;
	align-items: center;
	height: 100%;
}
</style>
