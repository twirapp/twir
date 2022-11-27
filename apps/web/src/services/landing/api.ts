export interface StatInfo {
  count: number;
  name: string;
}

export const getStats = async (): Promise<StatInfo[]> => {
  const res = await fetch('/api/v1/stats', { method: 'GET' });
  return res.json();
};
