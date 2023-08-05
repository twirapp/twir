import { useLocalStorage } from '@vueuse/core';
import { type LayoutItem } from 'grid-layout-plus';

const version = '1';

export type StatsItem = LayoutItem & { visible: boolean }

export const useStats = () => useLocalStorage<StatsItem[]>(`twirDashboardStatsPositions-v${version}`, [
	{
		x: 0,
		y: 0,
		w: 4,
		h: 1,
		i: 'Stats',
		minW: 4,
		minH: 1,
		visible: true,
	},
	{
		x: 4,
		y: 0,
		w: 2,
		h: 1,
		i: 'Uptime',
		visible: true,
	},
	{
		x: 6,
		y: 0,
		w: 2,
		h: 1,
		i: 'Viewers',
		visible: true,
	},
	{
		x: 8,
		y: 0,
		w: 2,
		h: 1,
		i: 'Followers',
		visible: true,
	},
	{
		x: 10,
		y: 0,
		w: 2,
		h: 1,
		i: 'Messages',
		visible: true,
	},
	{
		x: 0,
		y: 1,
		w: 2,
		h: 1,
		i: 'Used emotes',
		visible: true,
	},
	{
		x: 2,
		y: 5,
		w: 2,
		h: 1,
		i: 'Requested songs',
		visible: true,
	},
]);
