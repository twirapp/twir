import { createGlobalState, useLocalStorage } from '@vueuse/core'

export const useSidebarCollapseStore = createGlobalState(() => {
	const isCollapsed = useLocalStorage('twirSidebarIsCollapsed', false)

	function set(v: boolean) {
		isCollapsed.value = v
	}

	function toggle() {
		isCollapsed.value = !isCollapsed.value
	}

	return {
		isCollapsed,
		set,
		toggle,
	}
})
