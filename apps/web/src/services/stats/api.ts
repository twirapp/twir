export interface StatInfo {
  count: number
  name: string
}


export const getStats = async (): Promise<StatInfo[]> => {
  const req = await fetch('/api/v1/stats', { method: 'get' })
  return req.json();
};