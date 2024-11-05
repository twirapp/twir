<script setup lang="ts">
import {
	NConfigProvider,
	NDialogProvider,
	NMessageProvider,
	NNotificationProvider,
	darkTheme,
	lightTheme,
} from 'naive-ui'
import { computed, ref } from 'vue'
import { RouterView, useRouter } from 'vue-router'

import SidebarFloatingButton from './sidebar/sidebar-floating-button.vue'

import { Toaster as Sonner } from '@/components/ui/sonner'
import { Toaster } from '@/components/ui/toast'
import { useTheme } from '@/composables/use-theme.js'
import Sidebar from '@/layout/sidebar/sidebar.vue'

const { theme } = useTheme()
const themeStyles = computed(() => theme.value === 'dark' ? darkTheme : lightTheme)

const isRouterReady = ref(false)
const router = useRouter()

router.isReady().finally(() => isRouterReady.value = true)
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
					<Sidebar>
						<SidebarFloatingButton />
						<RouterView v-slot="{ Component, route }">
							<transition :name="route.meta.transition as string || 'router'" mode="out-in">
								<div
									:key="route.path"
									:style="{
										padding: route.meta?.noPadding ? undefined : '24px',
										height: '100%',
									}"
									class="dark:bg-[#101014]"
								>
									<component :is="Component" />
								</div>
							</transition>
						</RouterView>

						<Toaster />
						<Sonner />
					</Sidebar>
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
