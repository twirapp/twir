import { persistentAtom } from '@nanostores/persistent';
import { AuthUser } from '@tsuwari/shared';

type Dashboard = AuthUser['dashboards'][0];

export const selectedDashboardStore = persistentAtom<Dashboard>('selectedDashboard', null as any, {
  encode: JSON.stringify,
  decode: JSON.parse,
});

export const setSelectedDashboard = (value: Dashboard) => selectedDashboardStore.set(value);