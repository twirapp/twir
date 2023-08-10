import { useLocalStorage } from '@vueuse/core';
import { type LayoutItem } from 'grid-layout-plus';

const version = '2';

export type StatsItem = LayoutItem & { visible: boolean }

export const useStats = () => useLocalStorage<StatsItem[]>(`twirDashboardStatsPositions-v${version}`, [
	{ 'x': 0, 'y': 0, 'w': 4, 'h': 1, 'i': 'Stats', 'visible': true, 'moved': false },
	{ 'x': 4, 'y': 0, 'w': 1, 'h': 1, 'i': 'Uptime', 'visible': true, 'moved': false },
	{ 'x': 9, 'y': 0, 'w': 1, 'h': 1, 'i': 'Viewers', 'visible': true, 'moved': false },
	{ 'x': 5, 'y': 0, 'w': 1, 'h': 1, 'i': 'Followers', 'visible': true, 'moved': false },
	{ 'x': 6, 'y': 0, 'w': 1, 'h': 1, 'i': 'Messages', 'visible': true, 'moved': false },
	{ 'x': 8, 'y': 0, 'w': 1, 'h': 1, 'i': 'Used emotes', 'visible': true, 'moved': false },
	{ 'x': 10, 'y': 0, 'w': 2, 'h': 1, 'i': 'Requested songs', 'visible': true, 'moved': false },
	{ 'x': 7, 'y': 0, 'w': 1, 'h': 1, 'i': 'Subs', 'visible': true, 'moved': false },
]);
