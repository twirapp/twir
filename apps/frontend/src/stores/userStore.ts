import { persistentAtom } from '@nanostores/persistent';
import { AuthUser } from '@tsuwari/shared';
import { atom } from 'nanostores';

type Dashboard = AuthUser['dashboards'][0];

export const userStore = atom<(AuthUser & { apiKey: string }) | null | undefined>(null);
export const selectedDashboardStore = persistentAtom<Dashboard>('selectedDashboard', null as any, {
  encode: JSON.stringify,
  decode: JSON.parse,
});

export function setUser(data: AuthUser | null) {
  if (data?.id) {
    const dashboard = {
      id: '0',
      channelId: data.id,
      userId: data.id,
      twitchUser: { ...data, dashboards: undefined },
    };
    data.dashboards?.push(dashboard);

    if (!selectedDashboardStore.get()) {
      setSelectedDashboard(dashboard);
    }
  }

  userStore.set(data);
}

export const setSelectedDashboard = (value: Dashboard) => selectedDashboardStore.set(value);
