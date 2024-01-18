import { useLocalStorage } from '@vueuse/core';
import { type LayoutItem } from 'grid-layout-plus';

const version = '7';

export type WidgetItem = LayoutItem & { visible: boolean }

export const useWidgets = () => useLocalStorage<WidgetItem[]>(`twirWidgetsPositions-v${version}`, [
	{
		x: 6,
		y: 0,
		w: 6,
		h: 26,
		i: 'chat',
		minW: 3,
		minH: 8,
		visible: true,
	},
	{
		x: 0,
		y: 13,
		w: 6,
		h: 13,
		i: 'events',
		minW: 3,
		minH: 8,
		visible: true,
	},
	{
		x: 0,
		y: 0,
		w: 6,
		h: 13,
		i: 'stream',
		minW: 3,
		minH: 8,
		visible: true,
	},
]);
