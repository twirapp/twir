import { useLocalStorage } from '@vueuse/core';
import { type LayoutItem } from 'grid-layout-plus';

const version = '5';

export type WidgetItem = LayoutItem & { visible: boolean }

export const useWidgets = () => useLocalStorage<WidgetItem[]>(`twirWidgetsPositions-v${version}`, [
	{
		x: 0,
		y: 0,
		w: 4,
		h: 13,
		i: 'chat',
		minW: 3,
		minH: 8,
		visible: true,
	},
	{
		x: 4,
		y: 0,
		w: 4,
		h: 4,
		i: 'bot',
		minW: 3,
		minH: 4,
		visible: true,
	},
	{
		x: 8,
		y: 0,
		w: 4,
		h: 13,
		i: 'events',
		minW: 3,
		minH: 8,
		visible: true,
	},
	{
		x: 4,
		y: 4,
		w: 4,
		h: 9,
		i: 'stream',
		minW: 3,
		minH: 8,
		visible: true,
	},
]);
