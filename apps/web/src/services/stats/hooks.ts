import { useQuery } from '@tanstack/vue-query';

import { getStats } from './api.js';

export const useStats = () =>
  useQuery(['v1/query'], getStats, {
    retry: false,
    refetchOnReconnect: true,
    refetchOnWindowFocus: true,
    refetchInterval: 2500,
  });