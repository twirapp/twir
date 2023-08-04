import { useLocalStorage } from '@vueuse/core';

const version = '3';

export const usePositions = () => useLocalStorage(`twirWidgetsPositions-v${version}`, {
	chat: {
		x: 20,
		y: 20,
		width: 300,
		height: 400,
		isActive: false,
	},
	botManage: {
		x: 350,
		y: 20,
		width: 330,
		height: 200,
		isActive: false,
	},
	events: {
		x: 350,
		y: 250,
		width: 400,
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
