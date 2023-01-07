import { authorizeByTwitchCode } from './api.js';

export const handleTwitchLoginCallback = async () => {
  const params = new URLSearchParams(window.location.search);

  const error = params.get('error');
  if (error !== null) {
    throw new Error(`Got an error from twitch: "${error}"`);
  }

  const code = params.get('code');
  if (!code) throw new Error('Cannot find twitch code in query url');

  return await authorizeByTwitchCode(code);
};
