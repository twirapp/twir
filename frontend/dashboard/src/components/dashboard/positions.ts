import { useLocalStorage } from '@vueuse/core';

export const usePositions = () => useLocalStorage('twirWidgetsPositions', {
	chat: {
		x: 20,
		y: 20,
		width: 300,
		height: 400,
		isActive: false,
	},
	botManage: {
		x: 460,
		y: 20,
		width: 330,
		height: 200,
		isActive: false,
	},
});
