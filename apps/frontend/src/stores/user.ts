import { AuthUser } from '@tsuwari/shared';
import { atom } from 'nanostores';

import { setSelectedDashboard, selectedDashboardStore } from './selectedDashboard';

export const userStore = atom<AuthUser | null | undefined>(null);

export const setUser = (data: AuthUser | null) => {
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
};