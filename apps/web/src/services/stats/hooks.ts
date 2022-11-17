import { useQuery } from '@tanstack/vue-query';

import { getStats } from './api.js';

export const useStats = () =>
  useQuery(['stats'], getStats, {
    retry: false,
    refetchOnReconnect: true,
    refetchOnWindowFocus: true,
    refetchInterval: 2500,
  });