import { useQuery } from '@tanstack/react-query';
import { ChannelCustomvar } from '@twir/typeorm/entities/ChannelCustomvar';
import { getCookie } from 'cookies-next';
import { useContext } from 'react';

import { authFetcher } from '@/services/api/twirp.js';
import { SelectedDashboardContext } from '@/services/selectedDashboardProvider';


type Variable = {
  name: string;
  example?: string;
  description?: string
  visible: boolean
}

export const useVariables = () => {
  const dashboard = useContext(SelectedDashboardContext);
  const getUrl = () => `/api/v1/channels/${dashboard.id}/variables`;

  return useQuery<Variable[]>({
    queryKey: ['variablesList'],
    queryFn: async () => {
      const [custom, builtIn] = await Promise.all([
        authFetcher<ChannelCustomvar[]>(getUrl()),
        authFetcher<Variable[]>(`${getUrl()}/builtin`),
      ]);

      const list: Variable[] = [
        ...custom.map(v => ({
          name: v.name,
          example: `customvar|${v.name}`,
          description: `Your created variable ${v.name.toUpperCase()}`,
          visible: true,
        })),
        ...builtIn.filter(v => v.visible).sort((a, b) => a.name < b.name ? -1 : 1),
      ];

      return list;
    },
  });
};
