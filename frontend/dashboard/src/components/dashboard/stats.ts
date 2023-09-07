import { useLocalStorage } from '@vueuse/core';
import { type LayoutItem } from 'grid-layout-plus';

const version = '3';

export type StatsItem = LayoutItem & { visible: boolean }

export const useStats = () => useLocalStorage<StatsItem[]>(`twirDashboardStatsPositions-v${version}`, [
	{ 'x': 0, 'y': 0, 'w': 4, 'h': 1, 'i': 'streamInfo', 'visible': true, 'moved': false },
	{ 'x': 4, 'y': 0, 'w': 1, 'h': 1, 'i': 'uptime', 'visible': true, 'moved': false },
	{ 'x': 9, 'y': 0, 'w': 1, 'h': 1, 'i': 'viewers', 'visible': true, 'moved': false },
	{ 'x': 5, 'y': 0, 'w': 1, 'h': 1, 'i': 'followers', 'visible': true, 'moved': false },
	{ 'x': 6, 'y': 0, 'w': 1, 'h': 1, 'i': 'messages', 'visible': true, 'moved': false },
	{ 'x': 8, 'y': 0, 'w': 1, 'h': 1, 'i': 'usedEmotes', 'visible': true, 'moved': false },
	{ 'x': 10, 'y': 0, 'w': 2, 'h': 1, 'i': 'requestedSongs', 'visible': true, 'moved': false },
	{ 'x': 7, 'y': 0, 'w': 1, 'h': 1, 'i': 'subs', 'visible': true, 'moved': false },
]);
