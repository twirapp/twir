import { persistentAtom } from '@nanostores/persistent';
import { DashboardAccess } from '@tsuwari/prisma';
import { HelixUser } from '@twurple/api';
import { rawDataSymbol } from '@twurple/common';
import { atom } from 'nanostores';

type Dashboard = DashboardAccess & { twitch: HelixUser[typeof rawDataSymbol] }

export type User = HelixUser[typeof rawDataSymbol] & {
  dashboards: Array<Dashboard>
}

export const userStore = atom<User | null | undefined>(null);
export const selectedDashboardStore = persistentAtom<Dashboard>('selectedDashboard', null as any, {
  encode: JSON.stringify,
  decode: JSON.parse,
});

export function setUser(data: User | null) {
  if (data?.id) {
    const dashboard = { id: '0', channelId: data.id, userId: data.id, twitch: { ...data, dashboards: undefined } };
    data.dashboards?.push(dashboard);

    if (!selectedDashboardStore.get()) {
      setSelectedDashboard(dashboard);
    }
  }

  userStore.set(data);
}

export const setSelectedDashboard = (value: Dashboard) => selectedDashboardStore.set(value);