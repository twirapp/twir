<script setup lang="ts">
import { breakpointsTailwind, useBreakpoints } from '@vueuse/core'
import { ref, watch } from 'vue'
import { useRoute } from 'vue-router'

import HamburgerMenu from '@/components/ui/hamburger-menu.vue'
import Drawer from '@/layout/drawer'

const route = useRoute()
const isDrawerOpen = ref(false)

const breakPoints = useBreakpoints(breakpointsTailwind)
const isDesktopWindow = breakPoints.greaterOrEqual('md')

const handleCloseDrawer = () => isDrawerOpen.value = false
const handleToggleDrawer = () => isDrawerOpen.value = !isDrawerOpen.value

watch(isDesktopWindow, (v) => {
	if (!v && isDrawerOpen.value) {
		handleCloseDrawer()
	}
})

watch(() => route.fullPath, () => {
	if (!isDrawerOpen.value) return

	handleCloseDrawer()
})

// TODO: Don`t close if dropdown-profile-options is open
// useEventListener('keydown', (ev) => {
// if (ev.code !== 'Escape' || !isDrawerOpen.value) return
//
// handleCloseDrawer();
// })
</script>

<template>
	<HamburgerMenu :is-open="isDrawerOpen" @click="handleToggleDrawer" />

	<Drawer :show="isDrawerOpen" />
</template>
