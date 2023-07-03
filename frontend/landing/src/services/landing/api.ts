import { unprotectedApiClient } from '@/services/apiClients.js';

export interface StatInfo {
  count: number;
  name: string;
}

export const getStats = async (): Promise<StatInfo[]> => {
	const res = await unprotectedApiClient?.getStats({});

  return Object.entries(res.response).map(([name, value]) => ({ name, count: value }));
};
