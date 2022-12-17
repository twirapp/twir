export const swrFetcher = async (resource: RequestInfo | URL, init: RequestInit | undefined, isRetry = false): Promise<any> => {
  const req = await fetch(resource, {
    ...init,
    headers: {
      ...init?.headers,
      Authorization: 'window' in globalThis ? `Bearer ${localStorage.getItem('accessToken')}` : '',
    },
  });

  let data: any;
  if (req.headers.get('content-type') === 'application/json') {
    data = await req.json();
  } else {
    data = await req.text();
  }

  if (!req.ok && data.messages && data.messages[0].includes('invalid token') && !isRetry) {
    await refreshToken();
    return swrFetcher(resource, init);
  }

  return data;
};


const refreshToken = async () => {
  const req = await fetch('/api/auth/token', { method: 'POST' });
  if (!req.ok) {
    throw new Error('cannot refresh token');
  }

  const data = await req.json();
  localStorage.setItem('accessToken', data.accessToken);
};