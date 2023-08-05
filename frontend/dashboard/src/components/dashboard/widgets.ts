import { useLocalStorage } from '@vueuse/core';
import { type LayoutItem } from 'grid-layout-plus';

const version = '4';

export type WidgetItem = LayoutItem & { visible: boolean }

export const useWidgets = () => useLocalStorage<WidgetItem[]>(`twirWidgetsPositions-v${version}`, [
	{
		x: 0,
		y: 0,
		w: 4,
		h: 13,
		i: 'Chat',
		minW: 3,
		minH: 8,
		visible: true,
	},
	{
		x: 4,
		y: 0,
		w: 4,
		h: 4,
		i: 'Bot',
		minW: 3,
		minH: 4,
		visible: true,
	},
	{
		x: 8,
		y: 0,
		w: 4,
		h: 13,
		i: 'Events',
		minW: 3,
		minH: 8,
		visible: true,
	},
	{
		x: 4,
		y: 4,
		w: 4,
		h: 9,
		i: 'Stream',
		minW: 3,
		minH: 8,
		visible: true,
	},
]);
