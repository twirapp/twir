import { useLocalStorage, useSessionStorage } from '@mantine/hooks';
import type { Dashboard } from '@tsuwari/shared';
import { setCookie } from 'cookies-next';

// Local storage key
export const SELECTED_DASHBOARD_KEY = 'selectedDashboard';

export const useSelectedDashboard = () =>
   useLocalStorage<Dashboard | null>({
    key: SELECTED_DASHBOARD_KEY,
    serialize: (v) => JSON.stringify(v),
    deserialize: (v) => JSON.parse(v),
  });