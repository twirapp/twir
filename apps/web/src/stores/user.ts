import { persistentAtom } from '@nanostores/persistent';
import type { AuthUser, Dashboard } from '@tsuwari/shared';
import { atom } from 'nanostores';

import { createUserDashboard } from '@/services/dashboard.service.js';

export const userStore = atom<AuthUser | null>(null);

export const selectedDashboardStore = persistentAtom<Dashboard | null>('selected_dashboard', null, {
  encode: JSON.stringify,
  decode: JSON.parse,
});

export const accessTokenStore = persistentAtom('access_token');

export function setUser(user: AuthUser) {
  const userDashboard = createUserDashboard(user);
  user.dashboards.push(userDashboard);

  if (!selectedDashboardStore.get()) {
    selectedDashboardStore.set(userDashboard);
  }

  userStore.set(user);
}
