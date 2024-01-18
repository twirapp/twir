import { useLocalStorage } from '@vueuse/core';
import { defineStore } from 'pinia';

export const useSidebarCollapseStore = defineStore('isSidebarCollapsed', () => {
	const isCollapsed = useLocalStorage('twirSidebarIsCollapsed', false);

	function set(v: boolean) {
		isCollapsed.value = v;
	}

	function toggle() {
		isCollapsed.value = !isCollapsed.value;
	}

	return {
		isCollapsed,
		set,
		toggle,
	};
});
