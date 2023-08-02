import { useLocalStorage } from '@vueuse/core';

const version = '2';

export const usePositions = () => useLocalStorage(`twirWidgetsPositions-v${version}`, {
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
	stream: {
		x: 20,
		y: 500,
		width: 330,
		height: 200,
		isActive: false,
	},
});
