import { createContext, Dispatch, SetStateAction } from 'react';

export const SelectedDashboardContext = createContext({
} as {
  id: string,
  setId: Dispatch<SetStateAction<string>>
});