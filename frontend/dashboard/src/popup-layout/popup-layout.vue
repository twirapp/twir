<script setup lang="ts">
import {
	NConfigProvider,
	NDialogProvider,
	NMessageProvider,
	NNotificationProvider,
	darkTheme,
	lightTheme,
} from 'naive-ui'
import { computed } from 'vue'

import { Toaster as Sonner } from '@/components/ui/sonner'
import { useTheme } from '@/composables/use-theme'

const { theme } = useTheme()
const themeStyles = computed(() => (theme.value === 'dark' ? darkTheme : lightTheme))
</script>

<template>
	<NConfigProvider
		:theme="themeStyles"
		class="h-full"
		:breakpoints="{ xs: 0, s: 640, m: 1024, l: 1280, xl: 1536, xxl: 1920, '2xl': 2560 }"
	>
		<NNotificationProvider>
			<NMessageProvider :duration="2500" :closable="true">
				<NDialogProvider>
					<RouterView />
					<Sonner />
				</NDialogProvider>
			</NMessageProvider>
		</NNotificationProvider>
	</NConfigProvider>
</template>
