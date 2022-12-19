import { useLocalStorage } from '@mantine/hooks';
import type { Dashboard } from '@tsuwari/shared';

// Local storage key
const SELECTED_DASHBOARD_KEY = 'selectedDashboard';

export const useSelectedDashboard = () =>
  useLocalStorage<Dashboard | null>({
    key: SELECTED_DASHBOARD_KEY,
    serialize: (v) => JSON.stringify(v),
    deserialize: (v) => JSON.parse(v),
  });
